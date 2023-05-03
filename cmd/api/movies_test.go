package main

import (
	"encoding/json"
	"net/http"
	"testing"

	"greenlight.bcc/internal/assert"
)

func TestShowMovie(t *testing.T) {
	app := newTestApplication(t)

	ts := newTestServer(t, app.routesTest())
	defer ts.Close()

	tests := []struct {
		name     string
		urlPath  string
		wantCode int
		wantBody string
	}{
		{
			name:     "Valid ID",
			urlPath:  "/v1/movies/1",
			wantCode: http.StatusOK,
		},
		{
			name:     "Non-existent ID",
			urlPath:  "/v1/movies/2",
			wantCode: http.StatusNotFound,
		},
		{
			name:     "Negative ID",
			urlPath:  "/v1/movies/-1",
			wantCode: http.StatusNotFound,
		},
		{
			name:     "Decimal ID",
			urlPath:  "/v1/movies/1.23",
			wantCode: http.StatusNotFound,
		},
		{
			name:     "String ID",
			urlPath:  "/v1/movies/foo",
			wantCode: http.StatusNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			code, _, body := ts.get(t, tt.urlPath)

			assert.Equal(t, code, tt.wantCode)

			if tt.wantBody != "" {
				assert.StringContains(t, body, tt.wantBody)
			}

		})
	}

}

func TestListMovies(t *testing.T) {
	app := newTestApplication(t)
	ts := newTestServer(t, app.routesTest())
	defer ts.Close()

	tests := []struct {
		name     string
		urlPath  string
		wantCode int
		wantBody string
	}{
		{
			name:     "List Movies",
			urlPath:  "/v1/movies",
			wantCode: http.StatusOK,
			wantBody: `{"metadata":{"current_page":1,"page_size":20,"first_page":1,"last_page":1,"total_records":3},"movies":[{"id":1,"title":"Fight Club","year":1994,"runtime":"105 mins","genres":["Sigma"],"version":0},{"id":2,"title":"Drive","year":2014,"runtime":"120 mins","genres":["Ryan","Gosling"],"version":0},{"id":3,"title":"BladeRunner 2049","year":2012,"runtime":"179 mins","genres":["Oh, You don't even smile"],"version":0}]}`,
		},
		{
			name:     "List Movies",
			urlPath:  "/v1/movies",
			wantCode: http.StatusOK,
			wantBody: `{"metadata":{"current_page":1,"page_size":20,"first_page":1,"last_page":1,"total_records":3},"movies":[{"id":1,"title":"Fight Club","year":1994,"runtime":"105 mins","genres":["Sigma"],"version":0},{"id":2,"title":"Drive","year":2014,"runtime":"120 mins","genres":["Ryan","Gosling"],"version":0},{"id":3,"title":"BladeRunner 2049","year":2012,"runtime":"179 mins","genres":["Oh, You don't even smile"],"version":0}]}`,
		},
		{
			name:     "List Movies",
			urlPath:  "/v1/movies?sort=23",
			wantCode: http.StatusUnprocessableEntity,
			wantBody: "",
		},
		{
			name:     "List Movies",
			urlPath:  "/v1/movies?page=l",
			wantCode: http.StatusUnprocessableEntity,
			wantBody: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			code, _, body := ts.get(t, tt.urlPath)

			assert.Equal(t, code, tt.wantCode)

			if tt.wantBody != "" {
				assert.StringContains(t, body, tt.wantBody)
			}

		})
	}

}

func TestCreateMovie(t *testing.T) {
	app := newTestApplication(t)
	ts := newTestServer(t, app.routesTest())
	defer ts.Close()

	const (
		validTitle   = "Test Title"
		validYear    = 2021
		validRuntime = "105 mins"
	)

	validGenres := []string{"comedy", "drama"}

	tests := []struct {
		name     string
		Title    string
		Year     int32
		Runtime  string
		Genres   []string
		wantCode int
	}{
		{
			name:     "Valid submission",
			Title:    validTitle,
			Year:     validYear,
			Runtime:  validRuntime,
			Genres:   validGenres,
			wantCode: http.StatusCreated,
		},
		{
			name:     "Empty Title",
			Title:    "",
			Year:     validYear,
			Runtime:  validRuntime,
			Genres:   validGenres,
			wantCode: http.StatusUnprocessableEntity,
		},
		{
			name:     "year < 1888",
			Title:    validTitle,
			Year:     1500,
			Runtime:  validRuntime,
			Genres:   validGenres,
			wantCode: http.StatusUnprocessableEntity,
		},
		{
			name:     "test for wrong input",
			Title:    validTitle,
			Year:     validYear,
			Runtime:  validRuntime,
			Genres:   validGenres,
			wantCode: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			inputData := struct {
				Title   string   `json:"title"`
				Year    int32    `json:"year"`
				Runtime string   `json:"runtime"`
				Genres  []string `json:"genres"`
			}{
				Title:   tt.Title,
				Year:    tt.Year,
				Runtime: tt.Runtime,
				Genres:  tt.Genres,
			}

			b, err := json.Marshal(&inputData)
			if err != nil {
				t.Fatal("wrong input data")
			}
			if tt.name == "test for wrong input" {
				b = append(b, 'a')
			}

			code, _, _ := ts.postForm(t, "/v1/movies", b)

			assert.Equal(t, code, tt.wantCode)

		})
	}
}

func TestDeleteMovie(t *testing.T) {
	app := newTestApplication(t)
	ts := newTestServer(t, app.routesTest())
	defer ts.Close()

	tests := []struct {
		name     string
		urlPath  string
		wantCode int
		wantBody string
	}{
		{
			name:     "deleting existing movie",
			urlPath:  "/v1/movies/1",
			wantCode: http.StatusOK,
		},
		{
			name:     "Non-existent ID",
			urlPath:  "/v1/movies/2",
			wantCode: http.StatusNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			code, _, body := ts.deleteReq(t, tt.urlPath)

			assert.Equal(t, code, tt.wantCode)

			if tt.wantBody != "" {
				assert.StringContains(t, body, tt.wantBody)
			}

		})
	}

}

func TestUpdateMovieHandler(t *testing.T) {
	app := newTestApplication(t)
	ts := newTestServer(t, app.routesTest())
	defer ts.Close()

	// update movie fields
	tests := []struct {
		name     string
		url      string
		wantCode int
		Title    string   `json:"title"`
		Year     int32    `json:"year"`
		Runtime  string   `json:"runtime"`
		Genres   []string `json:"genres"`
	}{
		{
			name:     "Updated title",
			url:      "/v1/movies/1",
			wantCode: http.StatusOK,
			Title:    "Updated Title",
		},
		{
			name:     "Updated genres",
			url:      "/v1/movies/1",
			wantCode: http.StatusOK,
			Genres:   []string{"comedy, action"},
		},
		{
			name:     "Updated runtime",
			url:      "/v1/movies/1",
			wantCode: http.StatusOK,
			Runtime:  "140 mins",
		},
		{
			name:     "Invalid runtime format",
			url:      "/v1/movies/1",
			wantCode: http.StatusBadRequest,
			Runtime:  "140",
		},
		{
			name:     "Invalid year format",
			url:      "/v1/movies/1",
			wantCode: http.StatusUnprocessableEntity,
			Year:     234245,
		},
		{
			name:     "Movie not found",
			url:      "/v1/movies/2",
			wantCode: http.StatusNotFound,
			Runtime:  "140",
		},
		{
			name:     "Movie not found",
			url:      "/v1/movies/fff",
			wantCode: http.StatusNotFound,
			Runtime:  "140",
		},
		{
			name:     "Failed Validation",
			url:      "/v1/movies/2",
			wantCode: http.StatusNotFound,
			Runtime:  "18999 mins",
			Title:    "",
		},
		{
			name:     "Method Not Allowed test",
			url:      "/v1/movies/2",
			wantCode: http.StatusMethodNotAllowed,
			Runtime:  "183 min",
		}}

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
				Runtime: tt.Runtime,
				Genres:  tt.Genres,
			}

			b, err := json.Marshal(&inputData)
			if err != nil {
				t.Fatal("wrong input data")
			}

			if tt.name == "Method Not Allowed test" {
				code, _, _ := ts.put(t, tt.url, b)

				assert.Equal(t, code, tt.wantCode)
				return
			}

			code, _, _ := ts.patch(t, tt.url, b)

			assert.Equal(t, code, tt.wantCode)
		})
	}

}
