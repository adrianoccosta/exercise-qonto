package swaggerhdl

import (
	"fmt"
	"github.com/adrianoccosta/exercise-qonto/tools"
	"net/http"
	"os"
)

const (
	swaggerProcess = "swagger"
	swaggerEvent   = "get swagger"
)

// New returns a new swagger implementation
func New(rootDir string) Swagger {
	return swagger{
		rootDir: rootDir,
	}
}

// Swagger defines the swagger interface
type Swagger interface {
	ServeSwagger(w http.ResponseWriter, r *http.Request)
}

type swagger struct {
	rootDir string
}

func (s swagger) ServeSwagger(w http.ResponseWriter, r *http.Request) {
	filePath := fmt.Sprintf("%s/api.swagger.yaml", s.rootDir)
	_, err := os.Stat(filePath)
	if err != nil {
		tools.WriteError(w, http.StatusNotFound, err)
		return
	}
	http.ServeFile(w, r, filePath)
}
