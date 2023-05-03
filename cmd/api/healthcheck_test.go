package main

import (
	"net/http"
	"testing"

	"greenlight.bcc/internal/assert"
)

func TestHealthCheck(t *testing.T) {
	app := newTestApplication(t)
	ts := newTestServer(t, app.routesTest())
	defer ts.Close()

	tests := []struct {
		name     string
		urlPath  string
		wantCode int
	}{
		{
			name:     "health check",
			urlPath:  "/v1/healthcheck",
			wantCode: http.StatusOK,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			code, _, _ := ts.get(t, tt.urlPath)

			assert.Equal(t, code, tt.wantCode)
		})
	}
}
