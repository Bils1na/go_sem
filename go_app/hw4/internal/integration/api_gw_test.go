package integration

import (
	"context"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"gitlab.com/robotomize/gb-golang/homework/03-04-umanager/internal/env"
)

func TestRunMain(t *testing.T) {
	ctx := context.Background()

	e, err := env.Setup(ctx)
	if err != nil {
		t.Fatalf("env.Setup() error: %v", err)
	}

	httpServer := e.ApiGWHTTPServer

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	w := httptest.NewRecorder()

	httpServer.Handler.ServeHTTP(w, req)

	if w.Result().StatusCode != http.StatusOK {
		t.Fatalf("expected status 200, got %d", w.Result().StatusCode)
	}
}

func TestMain(m *testing.M) {
	os.Exit(m.Run())
}
