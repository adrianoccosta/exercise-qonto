package healthhdl

import (
	"github.com/adrianoccosta/exercise-qonto/tools"
	"net/http"
)

const (
	// StatusOK .
	StatusOK = "OK"
	// StatusNOK .
	StatusNOK = "NOK"

	healthEvent = "check health"
)

// Validator defines the function that health validators need to implement to participate in checking if the service
// should be marked as healthy or not.
type Validator func() Response

// Response is format of Response for the service
type Response struct {
	Name    string      `json:"name"`
	Status  string      `json:"status"`
	Message interface{} `json:"message,omitempty"`
}

// Health is a read model for status views.
type Health interface {
	Health(w http.ResponseWriter, r *http.Request)
}

type health struct {
	serviceName      string
	serviceVersion   string
	buildTime        string
	commitVersion    string
	pipelineNumber   string
	healthValidators []Validator
}

// New creates a new health service.
func New(serviceName, serviceVersion, buildTime, commitVersion, pipelineNumber string,
	healthValidators ...Validator) Health {
	return health{
		serviceName:      serviceName,
		serviceVersion:   serviceVersion,
		buildTime:        buildTime,
		commitVersion:    commitVersion,
		pipelineNumber:   pipelineNumber,
		healthValidators: healthValidators,
	}
}

// Health Show service health info
// @Summary Show service health info
// @ID read-health
// @Tags tools
// @Produce json
// @Success 200 {array} Response
// @Router /health [get]
func (h health) Health(w http.ResponseWriter, r *http.Request) {
	statusCode := http.StatusOK
	data := []Response{
		{Name: "Name", Status: StatusOK, Message: h.serviceName},
		{Name: "Version", Status: StatusOK, Message: h.serviceVersion},
		{Name: "Commit SHA", Status: StatusOK, Message: h.commitVersion},
		{Name: "Pipeline number", Status: StatusOK, Message: h.pipelineNumber},
		{Name: "GitTime", Status: StatusOK, Message: h.buildTime},
	}
	for _, v := range h.healthValidators {
		resp := v()
		data = append(data, resp)
		if resp.Status != StatusOK {
			statusCode = http.StatusInternalServerError
		}
	}

	tools.WriteJSON(w, statusCode, data)
}
