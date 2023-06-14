package api

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/nmramorov/gophemart/internal/db"
	"github.com/nmramorov/gophemart/internal/jobmanager"
	"github.com/nmramorov/gophemart/internal/mocks"
)

func TestCookiesMiddleware(t *testing.T) {
	cursor := &db.Cursor{DBInterface: mocks.NewMock()}
	ctx := context.Background()
	manager := jobmanager.NewJobmanager(cursor, "http://localhost:8081", &ctx)
	handler := NewHandler(cursor, manager)
	ts := httptest.NewServer(handler)

	defer ts.Close()

	request := httptest.NewRequest(http.MethodGet, "http://localhost:8080/api/user/balance", nil)
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, request)
	res := w.Result()
	defer res.Body.Close()
	assert.Equal(t, 401, res.StatusCode)
}
