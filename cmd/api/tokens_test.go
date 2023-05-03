package main

import (
	"encoding/json"
	"net/http"
	"testing"

	"greenlight.bcc/internal/assert"
)

func TestCreateToken(t *testing.T) {
	app := newTestApplication(t)
	ts := newTestServer(t, app.routesTest())
	defer ts.Close()

	tests := []struct {
		name     string
		Email    string
		Password string
		wantCode int
	}{
		{
			name:     "Create token",
			Email:    "test@test.com",
			Password: "1234567a",
			wantCode: http.StatusCreated,
		},
		{
			name:     "Unauthorized",
			Email:    "wads@wasd.com",
			Password: "1234553242sdfsfa",
			wantCode: http.StatusUnauthorized,
		},
		{
			name:     "Test for record not found",
			Email:    "notfound@wasd.com",
			Password: "1234553242sdfsfa",
			wantCode: http.StatusUnauthorized,
		},
		{
			name:     "Validation Fail",
			Email:    "wasd",
			Password: "wasd",
			wantCode: http.StatusUnprocessableEntity,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			inputData := struct {
				Email    string `json:"email"`
				Password string `json:"password"`
			}{
				Email:    tt.Email,
				Password: tt.Password,
			}

			b, err := json.Marshal(&inputData)
			if err != nil {
				t.Fatal("wrong input data")
			}
			if tt.name == "test for wrong input" {
				b = append(b, 'a')
			}

			code, _, _ := ts.postForm(t, "/v1/tokens/authentication", b)

			assert.Equal(t, code, tt.wantCode)
		})
	}

	code, _, _ := ts.postForm(t, "/v1/tokens/authentication", []byte{})

	assert.Equal(t, code, http.StatusBadRequest)
}
