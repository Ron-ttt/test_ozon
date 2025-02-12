package handlers

import (
	"bytes"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"testozon/internal/app/middleware"

	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_handlerWrapper_IndexPage(t *testing.T) {
	type want struct {
		code        int
		contentType string
	}
	tests := []struct {
		name    string
		request []byte
		want    want
	}{
		{
			name:    "positive test #1",
			request: []byte(`{"url":"http://localhost:8080/BpLnf"}`),
			want: want{
				code:        201,
				contentType: "application/json",
			},
		},
		{
			name:    "negative test #1",
			request: []byte(""),
			want: want{
				code:        400,
				contentType: "text/plain; charset=utf-8",
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			hw := MInit()

			r := mux.NewRouter()
			r.Use(middleware.Logger1, middleware.GzipMiddleware)

			r.HandleFunc("/", hw.IndexPage)
			request := httptest.NewRequest(http.MethodPost, hw.baseURL, bytes.NewReader(test.request))
			w := httptest.NewRecorder()

			r.ServeHTTP(w, request)
			res := w.Result()
			defer res.Body.Close()

			assert.Equal(t, test.want.code, res.StatusCode)
			defer w.Result().Body.Close()

			resBody, err := io.ReadAll(res.Body)

			require.NoError(t, err)
			assert.NotEmpty(t, string(resBody))
			assert.Equal(t, test.want.contentType, res.Header.Get("Content-Type"))
		})
	}
}

func Test_handlerWrapper_Redirect(t *testing.T) {
	type want struct {
		code        int
		location    string
		contentType string
	}

	tests := []struct {
		name string
		id   string
		want want
	}{
		{
			name: "positive test #1",
			id:   "123456",
			want: want{
				code:        http.StatusTemporaryRedirect,
				location:    "http://love_nika",
				contentType: "",
			},
		},
		{
			name: "negative test #1",
			id:   "invalid",
			want: want{
				code:        http.StatusBadRequest,
				location:    "",
				contentType: "",
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			handler := MInit()

			r := mux.NewRouter()
			r.Use(middleware.Logger1, middleware.GzipMiddleware)

			r.HandleFunc("/{id}", handler.Redirect)
			w2 := httptest.NewRecorder()
			r.ServeHTTP(w2, httptest.NewRequest(http.MethodGet, handler.baseURL+test.id, nil))

			res := w2.Result()
			defer res.Body.Close()
			assert.Equal(t, test.want.code, res.StatusCode)
			defer w2.Result().Body.Close()

			location := w2.Header().Get("Location")
			assert.Equal(t, test.want.location, location)
		})
	}
}
