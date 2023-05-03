package main

import (
	"encoding/json"
	"net/http"
	"testing"

	"greenlight.bcc/internal/assert"
)

func TestRegisterUser(t *testing.T) {
	app := newTestApplication(t)
	ts := newTestServer(t, app.routesTest())
	defer ts.Close()

	tests := []struct {
		name     string
		Name     string
		Email    string
		Password string
		wantCode int
	}{
		{
			name:     "Register user",
			Name:     "Test",
			Email:    "anotherTest@test.com",
			Password: "TestTest123",
			wantCode: http.StatusCreated,
		},
		{
			name:     "Duplicate email",
			Name:     "Test",
			Email:    "test@test.com",
			Password: "TestTest123",
			wantCode: http.StatusUnprocessableEntity,
		},
		{
			name:     "test for wrong input",
			Name:     "Test",
			Email:    "test@test.com",
			Password: "TestTest123",
			wantCode: http.StatusBadRequest,
		},
		{
			name:     "Fail validation",
			Name:     "",
			Email:    "",
			Password: "",
			wantCode: http.StatusUnprocessableEntity,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			inputData := struct {
				Name     string `json:"name"`
				Email    string `json:"email"`
				Password string `json:"password"`
			}{
				Name:     tt.Name,
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

			code, _, _ := ts.postForm(t, "/v1/users", b)

			assert.Equal(t, code, tt.wantCode)

		})
	}
}

func TestActivateUser(t *testing.T) {
	app := newTestApplication(t)
	ts := newTestServer(t, app.routesTest())
	defer ts.Close()

	tests := []struct {
		name           string
		TokenPlainText string
		wantCode       int
	}{
		{
			name:           "Activate user",
			TokenPlainText: "BusinessManBusinessPlan123",
			wantCode:       http.StatusOK,
		},
		{
			name:           "Fail validation",
			TokenPlainText: "",
			wantCode:       http.StatusUnprocessableEntity,
		},

		{
			name:           "Test for record not found",
			TokenPlainText: "BusinessManBusinessPlan404",
			wantCode:       http.StatusUnprocessableEntity,
		},
		{
			name:           "test for wrong input",
			TokenPlainText: "",
			wantCode:       http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			inputData := struct {
				TokenPlainText string `json:"token"`
			}{
				TokenPlainText: tt.TokenPlainText,
			}

			b, err := json.Marshal(&inputData)
			if err != nil {
				t.Fatal("wrong input data")
			}
			if tt.name == "test for wrong input" {
				b = append(b, 'a')
			}

			code, _, _ := ts.put(t, "/v1/users/activated", b)

			assert.Equal(t, code, tt.wantCode)

		})
	}
}
