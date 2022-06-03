package cmd

import (
	"context"
	"fmt"
	"github.com/adrianoccosta/exercise-qonto/cmd/config"
	"github.com/adrianoccosta/exercise-qonto/internal/handlers/bankaccounthdl"
	"github.com/adrianoccosta/exercise-qonto/internal/handlers/healthhdl"
	"github.com/adrianoccosta/exercise-qonto/internal/handlers/metricshdl"
	"github.com/adrianoccosta/exercise-qonto/internal/handlers/swaggerhdl"
	"github.com/adrianoccosta/exercise-qonto/internal/handlers/transactionhdl"
	"github.com/adrianoccosta/exercise-qonto/internal/handlers/transferhdl"
	"github.com/adrianoccosta/exercise-qonto/internal/repository/bankaccountrepo"
	"github.com/adrianoccosta/exercise-qonto/internal/repository/transactionrepo"
	"github.com/adrianoccosta/exercise-qonto/internal/services/bankaccountsvc"
	"github.com/adrianoccosta/exercise-qonto/internal/services/transactionsvc"
	"github.com/adrianoccosta/exercise-qonto/internal/services/transfersvc"
	"github.com/adrianoccosta/exercise-qonto/log"
	"github.com/adrianoccosta/exercise-qonto/tools"
	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/urfave/cli/v2"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

const (
	timeout     = 30 * time.Second
	idleTimeout = 61 * time.Second

	listenPortProp    = "port"
	listenAddressProp = "addr"

	databaseFilePathProp     = "database-file-path"
	databaseMaxIdleConnsProp = "database-max-idle-conns"
	databaseMaxOpenConnsProp = "database-max-open-conns"
	databaseconnMaxLifetime  = "database-max-conn-lifetime"
)

// APICommand is the command to run the web server
var APICommand = &cli.Command{
	Name:   "api",
	Usage:  "service-qonto API",
	Action: runAPICommand,
	Flags: []cli.Flag{
		&cli.IntFlag{Name: listenPortProp, Value: tools.EnvIntOrDefault("PORT", 8080), Usage: "listen port"},
		&cli.StringFlag{Name: listenAddressProp, Value: "0.0.0.0", Usage: "HTTP listen address"},
		&cli.StringFlag{Name: databaseFilePathProp, Value: tools.GetEnv("DATABASE_FILE_PATH"), Usage: "database file path (e.g., qonto.db)"},
		&cli.IntFlag{Name: databaseMaxIdleConnsProp, Value: tools.EnvIntOrDefault("DATABASE_MAX_IDLE_CONNS", 15), Usage: "database max idle connections (e.g., 15)"},
		&cli.IntFlag{Name: databaseMaxOpenConnsProp, Value: tools.EnvIntOrDefault("DATABASE_MAX_OPEN_CONNS", 15), Usage: "database max open connections (e.g., 15)"},
		&cli.IntFlag{Name: databaseconnMaxLifetime, Value: tools.EnvIntOrDefault("DATABASE_MAX_CONN_LIFETIME", 30), Usage: "database max connection lifetime in minutes (e.g., 5)"},
	},
}

func runAPICommand(ctx *cli.Context) error {

	addr := fmt.Sprintf("%s:%d", ctx.String(listenAddressProp), ctx.Int(listenPortProp))
	logger := ctx.App.Metadata["Logger"].(log.Logger)

	r := configAPIHandlers(ctx, logger)

	// server
	s := &http.Server{
		Addr:           addr,
		Handler:        r,
		ReadTimeout:    timeout,
		WriteTimeout:   timeout,
		IdleTimeout:    idleTimeout,
		MaxHeaderBytes: 1 << 20,
	}

	// Graceful kill
	serverErrors := make(chan error, 1)
	// Start the listener.
	go func() {
		logger.Info(fmt.Sprintf("HTTP is running: %v", addr))
		serverErrors <- s.ListenAndServe()
	}()

	// Make a channel to listen for an interrupt or terminate signal from the OS.
	// Use a buffered channel because the signal package requires it.
	osSignals := make(chan os.Signal, 1)
	signal.Notify(osSignals, os.Interrupt, syscall.SIGTERM)

	// Block waiting for a receive on either channel
	select {
	case err := <-serverErrors:
		logger.WithError(err).Fatal("error starting server")

	case <-osSignals:
		// Create a context to attempt a graceful 5 second shutdown.
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		// Attempt the graceful shutdown by closing the listener and
		// completing all inflight requests.
		if err := s.Shutdown(ctx); err != nil {
			logger.WithError(err).Info("Could not stop server gracefully")
			logger.Info("Initiating hard shutdown")
			if err := s.Close(); err != nil {
				logger.WithError(err).Fatal("could not stop http server")
			}
		}
	}

	logger.Info("Shut down successful")

	return nil
}

func configAPIHandlers(ctx *cli.Context, logger log.Logger) *mux.Router {

	buildTime := fmt.Sprint(ctx.App.Metadata["BuildTime"])
	commitVersion := fmt.Sprint(ctx.App.Metadata["CommitVersion"])
	pipelineNumber := fmt.Sprint(ctx.App.Metadata["PipelineNumber"])

	//process := fmt.Sprint(ctx.App.Metadata["Process"])
	prometheusNamespace := fmt.Sprint(ctx.App.Metadata["PrometheusNamespace"])
	prometheusSubsystem := fmt.Sprint(ctx.App.Metadata["PrometheusSubsystem"])

	r := mux.NewRouter()
	r.NotFoundHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	})
	r.MethodNotAllowedHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusMethodNotAllowed)
	})

	// prometheus integration
	// simpler/default metrics configuration
	metrics := metricshdl.NewMetrics(prometheusNamespace, prometheusSubsystem)

	// rds connection. Probably we need to set this to a common path
	conn := config.DBConnection{
		Path:              ctx.String(databaseFilePathProp),
		MaxIdleConns:      ctx.Int(databaseMaxIdleConnsProp),
		MaxOpenConns:      ctx.Int(databaseMaxOpenConnsProp),
		ConnMaxTTLMinutes: ctx.Int(databaseconnMaxLifetime),
	}
	rds := config.InitDBConnection(conn, logger)

	// Repository
	bankAccountRepository := bankaccountrepo.New(rds)
	transactionRepository := transactionrepo.New(rds)

	// services
	bankAccountService := bankaccountsvc.New(bankAccountRepository, logger)
	transactionService := transactionsvc.New(transactionRepository, logger)
	transferService := transfersvc.New(transactionRepository, bankAccountRepository, logger)

	// handlers
	handlerHealth := healthhdl.New(ctx.App.Name, ctx.App.Version, buildTime, commitVersion, pipelineNumber)
	handlerSwagger := swaggerhdl.New("./")

	handlerBankAccount := bankaccounthdl.New(bankAccountService, logger)
	handlertransaction := transactionhdl.New(transactionService, logger)
	handlerTransfer := transferhdl.New(transferService, logger)

	apiRouter := r.PathPrefix("/qonto/api").Subrouter()

	apiRouter.Handle("/metrics", promhttp.HandlerFor(metrics.Registry(), promhttp.HandlerOpts{}))
	apiRouter.HandleFunc("/health", handlerHealth.Health)
	apiRouter.HandleFunc("/swagger", handlerSwagger.ServeSwagger)

	apiV1Router := apiRouter.PathPrefix("/v1").Subrouter()

	//apiV1Router.Use(middleware.Logger(false, logger))
	apiV1Router.Use(metrics.Handler)
	apiV1Router.Use(func(h http.Handler) http.Handler {
		return http.TimeoutHandler(h, 5*time.Second, "timeout while invoking service")
	})

	handlerBankAccount.Handlers(apiV1Router)
	handlertransaction.Handlers(apiV1Router)
	handlerTransfer.Handlers(apiV1Router)

	return r
}
