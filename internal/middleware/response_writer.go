package middleware

import (
	"bytes"
	"net/http"
)

// responseWriter is a wrapper for http.ResponseWriter that allows the
// written HTTP statusCode code and body to be captured.
type responseWriter struct {
	http.ResponseWriter
	statusCode  int
	wroteHeader bool
	body        *bytes.Buffer
}

// WrapResponseWriter creates a new responseWriter that is able to capture HTTP statusCode and body.
// body capturing is optional and it will be captured only if a non nil buffer is passed to the method.
func WrapResponseWriter(w http.ResponseWriter, body *bytes.Buffer) *responseWriter {
	if body != nil {
		return &responseWriter{ResponseWriter: w, body: body}
	}

	return &responseWriter{ResponseWriter: w}
}

func (rw *responseWriter) Status() int {
	return rw.statusCode
}

func (rw *responseWriter) WriteHeader(code int) {
	if rw.wroteHeader {
		return
	}

	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
	rw.wroteHeader = true
}

func (rw *responseWriter) Write(b []byte) (int, error) {
	if rw.body != nil {
		rw.body.Write(b)
	}
	return rw.ResponseWriter.Write(b)
}
