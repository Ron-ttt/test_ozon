package handlers

import (
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"testozon/internal/app/middleware"

	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_handlerWrapper_IndexPage(t *testing.T) { // работает удивительно ведь не я это делала
	type want struct { // я не ебу что для негативного надо работает и заебись
		code        int
		request     string
		contentType string
	}
	tests := []struct {
		name string
		want want
	}{
		{
			name: "positive test #1",
			want: want{
				code:        201,
				request:     "https://example.com",
				contentType: "text/plain",
			},
		},
		{
			name: "negative test #1",
			want: want{
				code:        400,
				request:     "",
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
			w2 := strings.NewReader(test.want.request)
			request := httptest.NewRequest(http.MethodPost, hw.baseURL, w2)
			// создаём новый Recorder
			w := httptest.NewRecorder()
			r.ServeHTTP(w, request)
			res := w.Result()
			defer res.Body.Close()
			// проверяем код ответа
			assert.Equal(t, test.want.code, res.StatusCode)
			// получаем и проверяем тело запроса
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
			// Создаем тестовый обработчик
			handler := MInit()

			// Создаем тестовый запрос

			r := mux.NewRouter()
			r.Use(middleware.Logger1, middleware.GzipMiddleware)

			r.HandleFunc("/{id}", handler.Redirect)
			// Выполняем запрос
			w2 := httptest.NewRecorder()
			r.ServeHTTP(w2, httptest.NewRequest(http.MethodGet, handler.baseURL+test.id, nil))

			// Проверяем код ответа
			res := w2.Result()
			defer res.Body.Close()
			assert.Equal(t, test.want.code, res.StatusCode)
			defer w2.Result().Body.Close()
			// Проверяем заголовок Location
			location := w2.Header().Get("Location")
			assert.Equal(t, test.want.location, location)
		})
	}
}
