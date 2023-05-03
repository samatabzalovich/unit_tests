package main

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestServerErrorResponse(t *testing.T) {
	app := newTestApplication(t)

	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodGet, "/", nil)

	err := errors.New("test error")
	app.serverErrorResponse(w, r, err)

	// check response code
	if w.Code != http.StatusInternalServerError {
		t.Errorf("expected response code %d; got %d", http.StatusInternalServerError, w.Code)
	}

	// check response headers
	expectedContentType := "application/json"
	if contentType := w.Header().Get("Content-Type"); contentType != expectedContentType {
		t.Errorf("expected Content-Type header %q; got %q", expectedContentType, contentType)
	}

	// check response body
	expectedBody := "{\"error\":\"the server encountered a problem and could not process your request\"}\n"
	if w.Body.String() != expectedBody {
		t.Errorf("expected response body %q; got %q", expectedBody, w.Body.String())
	}
}
