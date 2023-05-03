package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strings"
	"testing"

	"greenlight.bcc/internal/assert"
)

func TestRecoverPanicMiddleware(t *testing.T) {

	app := newTestApplication(t)

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	recorder := httptest.NewRecorder()

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		panic("something went wrong")
	})
	middleware := app.recoverPanic(handler)

	middleware.ServeHTTP(recorder, req)

	assert.Equal(t, http.StatusInternalServerError, recorder.Code)

	expectedJSON := `{"error":"the server encountered a problem and could not process your request"}`
	actualJSON := strings.TrimSpace(recorder.Body.String())

	if !json.Valid([]byte(actualJSON)) {
		t.Fatalf("invalid JSON response: %s", actualJSON)
	}

	var expected interface{}
	var actual interface{}
	if err := json.Unmarshal([]byte(expectedJSON), &expected); err != nil {
		t.Fatal(err)
	}
	if err := json.Unmarshal([]byte(actualJSON), &actual); err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(actual, expected) {
		t.Errorf("unexpected response body:\nexpected: %v\nactual: %v", expected, actual)
	}
}

func TestRateLimitMiddleware(t *testing.T) {
	app := newTestApplication(t)

	ts := newTestServer(t, app.routesTest())
	defer ts.Close()

	for i := 0; i < 4; i++ {

		t.Run("Valid", func(t *testing.T) {
			code, _, _ := ts.get(t, "/v1/movies")

			assert.Equal(t, code, http.StatusOK)
		})
	}

	code, _, _ := ts.get(t, "/v1/movies")
	assert.Equal(t, code, http.StatusTooManyRequests)
}

func TestRequireAuthenticated(t *testing.T) {
	app := newTestApplication(t)
	ts := newTestServer(t, app.routesTest())
	defer ts.Close()

	tests := []struct {
		name     string
		wantCode int
		token    string
	}{
		{
			name:     "OK test",
			wantCode: http.StatusOK,
			token:    "Bearer BusinessManBusinessPlan123",
		},
		{
			name:     "Unauthorized",
			wantCode: http.StatusUnauthorized,
			token:    "Bearer f",
		},
		{
			name:     "Anonymous user test",
			wantCode: http.StatusUnauthorized,
			token:    "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			code, _, _ := ts.getWithAuth(t, "/testauth/v1/movies", tt.token)
			assert.Equal(t, code, tt.wantCode)
		})
	}
}

func TestRequireActivated(t *testing.T) {
	app := newTestApplication(t)
	ts := newTestServer(t, app.routesTest())
	defer ts.Close()

	tests := []struct {
		name     string
		wantCode int
		token    string
	}{
		{
			name:     "OK test",
			wantCode: http.StatusOK,
			token:    "Bearer BusinessManBusinessPlan123",
		},
		{
			name:     "Unauthorized",
			wantCode: http.StatusUnauthorized,
			token:    "Bearer f",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			code, _, _ := ts.getWithAuth(t, "/testactivated/v1/movies", tt.token)
			assert.Equal(t, code, tt.wantCode)
		})
	}
}

func TestRequirePermissions(t *testing.T) {
	app := newTestApplication(t)
	ts := newTestServer(t, app.routesTest())
	defer ts.Close()

	tests := []struct {
		name     string
		wantCode int
		token    string
	}{
		{
			name:     "OK test movies:read",
			wantCode: http.StatusOK,
			token:    "Bearer BusinessManBusinessPlan123",
		},
		{
			name:     "Forbidden request (not activated account)",
			wantCode: http.StatusForbidden,
			token:    "Bearer BusinessManBusinessPlanNOO",
		},
		{
			name:     "Forbidden request (no required permissions)",
			wantCode: http.StatusForbidden,
			token:    "Bearer BusinessManBusinessPlan000",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			code, _, _ := ts.getWithAuth(t, "/testpermissions/v1/movies", tt.token)
			assert.Equal(t, code, tt.wantCode)
		})
	}
}

func TestAuthenticateMiddleware(t *testing.T) {
	app := newTestApplication(t)
	ts := newTestServer(t, app.routesTest())
	defer ts.Close()

	tests := []struct {
		name     string
		url      string
		wantCode int
		token    string
		Title    string   `json:"title"`
		Year     int32    `json:"year"`
		Genres   []string `json:"genres"`
		Runtime  string   `json:"runtime"`
	}{
		{
			name:     "Anonym",
			url:      "/v1/movies/1",
			wantCode: http.StatusOK,
			token:    "",
			Title:    "Updated Title",
			Runtime:  "105 mins",
		},
		{
			name:     "No Prefix Token",
			url:      "/v1/movies/1",
			wantCode: http.StatusUnauthorized,
			token:    "wasd",
			Title:    "Updated Title",
			Runtime:  "105 mins",
		},
		{
			name:     "Invalid Token",
			url:      "/v1/movies/1",
			wantCode: http.StatusUnauthorized,
			token:    "Bearer wasd",
			Title:    "Updated Title",
			Runtime:  "105 mins",
		}, {
			name:     "OK Token",
			url:      "/v1/movies/1",
			wantCode: http.StatusOK,
			token:    "Bearer BusinessManBusinessPlan123",
			Title:    "Updated Title",
			Runtime:  "105 mins",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			inputData := struct {
				Title   string   `json:"title,omitempty"`
				Year    int32    `json:"year,omitempty"`
				Runtime string   `json:"runtime,omitempty"`
				Genres  []string `json:"genres,omitempty"`
			}{
				Title:   tt.Title,
				Year:    tt.Year,
				Genres:  tt.Genres,
				Runtime: tt.Runtime,
			}

			b, err := json.Marshal(&inputData)
			if err != nil {
				t.Fatal("wrong input data")
			}

			code, _, _ := ts.patchForAuth(t, tt.url, b, tt.token)

			assert.Equal(t, code, tt.wantCode)
		})
	}
}
