package healthhdl

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

const (
	serviceName    = "service-qonto"
	serviceVersion = "version"
	buildTime      = "build"
	commitVersion  = "version1"
	pipelineNumber = "1"
)

var (
	testExpectedMasterOK = func() Response {
		return Response{
			Name:    "master",
			Status:  StatusOK,
			Message: "bla",
		}
	}

	testExpectedMasterNOK = func() Response {
		return Response{
			Name:    "master",
			Status:  StatusNOK,
			Message: "bla",
		}
	}
)

func TestHealth(t *testing.T) {
	req, err := http.NewRequest("GET", "/health", nil)
	if err != nil {
		t.Fatal(err)
	}

	h := New(serviceName, serviceVersion, buildTime, commitVersion, pipelineNumber, testExpectedMasterOK)

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(h.Health)

	handler.ServeHTTP(rr, req)

	expectedStatusCode := http.StatusOK
	if status := rr.Code; status != expectedStatusCode {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, expectedStatusCode)
	}

	var actual []Response
	err = json.NewDecoder(rr.Body).Decode(&actual)
	if err != nil {
		t.Fatal(err)
	}

	assert.EqualValues(t, "Name", actual[0].Name)
	assert.EqualValues(t, StatusOK, actual[0].Status)
	assert.EqualValues(t, "service-qonto", actual[0].Message)
}
