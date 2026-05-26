package transport_test

import (
	"bytes"
	"context"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/cubelitblade/community-v2/backend/services/authz/internal/authz"
	"github.com/cubelitblade/community-v2/backend/services/authz/internal/authz/transport"
	"github.com/gin-gonic/gin"
)

var testLogger = slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelError}))

type stubChecker struct {
	allowed bool
	err     error
}

func (c *stubChecker) Check(_ context.Context, _, _, _ string) (bool, error) {
	return c.allowed, c.err
}

func setupRouter(h *transport.Handler) *gin.Engine {
	gin.SetMode(gin.TestMode)

	router := gin.New()
	h.RegisterRoutes(router)

	return router
}

func doRequest(t *testing.T, router *gin.Engine, body string) *httptest.ResponseRecorder {
	t.Helper()

	var req *http.Request
	if body == "" {
		req = httptest.NewRequestWithContext(context.Background(), http.MethodPost, "/authz/check", http.NoBody)
	} else {
		req = httptest.NewRequestWithContext(context.Background(), http.MethodPost, "/authz/check", bytes.NewBufferString(body))
	}
	req.Header.Set("Content-Type", "application/json")

	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	return rec
}

func TestHandler_CheckAllowed(t *testing.T) {
	t.Parallel()

	h := transport.NewHandler(transport.Deps{
		Checker: &stubChecker{allowed: true, err: nil},
		Logger:  testLogger,
	})

	router := setupRouter(h)
	rec := doRequest(t, router, `{"user":"user:alice","relation":"can_edit","object":"post:1"}`)

	if rec.Code != http.StatusOK {
		t.Errorf("status = %d, want %d", rec.Code, http.StatusOK)
	}

	if body := rec.Body.String(); body != `{"allowed":true}` {
		t.Errorf("body = %s, want {\"allowed\":true}", body)
	}
}

func TestHandler_CheckDenied(t *testing.T) {
	t.Parallel()

	h := transport.NewHandler(transport.Deps{
		Checker: &stubChecker{allowed: false, err: nil},
		Logger:  testLogger,
	})

	router := setupRouter(h)
	rec := doRequest(t, router, `{"user":"user:bob","relation":"can_edit","object":"post:1"}`)

	if rec.Code != http.StatusOK {
		t.Errorf("status = %d, want %d", rec.Code, http.StatusOK)
	}

	if body := rec.Body.String(); body != `{"allowed":false}` {
		t.Errorf("body = %s, want {\"allowed\":false}", body)
	}
}

func TestHandler_CheckInvalidBody(t *testing.T) {
	t.Parallel()

	h := transport.NewHandler(transport.Deps{
		Checker: &stubChecker{allowed: false, err: nil},
		Logger:  testLogger,
	})

	router := setupRouter(h)
	rec := doRequest(t, router, `{"not":"valid"}`)

	if rec.Code != http.StatusBadRequest {
		t.Errorf("status = %d, want %d", rec.Code, http.StatusBadRequest)
	}
}

func TestHandler_CheckError(t *testing.T) {
	t.Parallel()

	h := transport.NewHandler(transport.Deps{
		Checker: &stubChecker{allowed: false, err: authz.ErrCheckFailed},
		Logger:  testLogger,
	})

	router := setupRouter(h)
	rec := doRequest(t, router, `{"user":"user:alice","relation":"can_edit","object":"post:1"}`)

	if rec.Code != http.StatusInternalServerError {
		t.Errorf("status = %d, want %d", rec.Code, http.StatusInternalServerError)
	}
}

func TestCheckRoute(t *testing.T) {
	t.Parallel()

	h := transport.NewHandler(transport.Deps{
		Checker: &stubChecker{allowed: false, err: nil},
		Logger:  testLogger,
	})

	router := setupRouter(h)
	rec := doRequest(t, router, "")

	if rec.Code != http.StatusBadRequest {
		t.Errorf("nil body: status = %d, want %d", rec.Code, http.StatusBadRequest)
	}
}
