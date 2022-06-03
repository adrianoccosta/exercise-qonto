package main

import (
	"github.com/adrianoccosta/exercise-qonto/cmd"
	"github.com/adrianoccosta/exercise-qonto/log"
	"github.com/urfave/cli/v2"
	"go.uber.org/zap"
	"os"
	"runtime"
)

const (
	// PrometheusNamespace define the value for the group app
	PrometheusNamespace = "exercise"
	// PrometheusSubsystem define the value for the sub group app
	PrometheusSubsystem = "qonto"
	// Process defines the value that should be used for this project
	Process = "qonto"
	// App export name of service
	App = "service-qonto"
)

var (
	// Version export version of service
	Version = "1.0.0"
	// BuildTime represents the build time
	BuildTime string
	// CommitVersion represents the last commit hash of the build
	CommitVersion string
	// PipelineNumber represents the last commit hash of the build
	PipelineNumber string
)

func init() {
	// sets the maximum number of CPUs that can be executing
	runtime.GOMAXPROCS(runtime.NumCPU())
}

func main() {
	app := cli.App{
		Name:        App,
		Version:     Version,
		Description: "service-qonto is the implementation exercise of the candidate Adriano Costa",
		Before:      setupBefore,
		Commands: []*cli.Command{
			cmd.APICommand,
		},
	}
	app.Flags = append(app.Flags, []cli.Flag{}...)
	if len(os.Args) == 1 {
		os.Args = append(os.Args, cmd.APICommand.Name)
	}
	if err := app.Run(os.Args); err != nil {
		panic(err)
	}
	os.Exit(0)
}

func setupBefore(cli *cli.Context) error {

	logger, err := log.New()
	if err != nil {
		panic("could not initialise log due to: " + err.Error())
	}

	fields := []zap.Field{
		zap.Any("service", cli.App.Name),
		zap.Any("version", cli.App.Version),
		zap.Any("commit_hash", CommitVersion),
		zap.Any("build_time", BuildTime),
	}

	logger = logger.With(fields...)

	cli.App.Metadata["Logger"] = logger
	cli.App.Metadata["Process"] = Process
	cli.App.Metadata["PrometheusNamespace"] = PrometheusNamespace
	cli.App.Metadata["PrometheusSubsystem"] = PrometheusSubsystem
	cli.App.Metadata["BuildTime"] = BuildTime
	cli.App.Metadata["CommitVersion"] = CommitVersion
	cli.App.Metadata["PipelineNumber"] = PipelineNumber

	return nil
}
