package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestBearerAuthAllowsValidKey(t *testing.T) {
	auth := NewBearerAuth([]string{"sk-valid"})
	called := false
	handler := auth.Wrap(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		called = true
		w.WriteHeader(http.StatusNoContent)
	}))
	req := httptest.NewRequest(http.MethodGet, "/api/v1/ping", nil)
	req.Header.Set("Authorization", "Bearer sk-valid")
	rr := httptest.NewRecorder()

	handler.ServeHTTP(rr, req)

	if !called {
		t.Fatal("wrapped handler was not called")
	}
	if rr.Code != http.StatusNoContent {
		t.Fatalf("status = %d, want %d", rr.Code, http.StatusNoContent)
	}
}

func TestBearerAuthRejectsMissingKey(t *testing.T) {
	auth := NewBearerAuth([]string{"sk-valid"})
	rr := httptest.NewRecorder()

	auth.Wrap(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t.Fatal("wrapped handler should not be called")
	})).ServeHTTP(rr, httptest.NewRequest(http.MethodGet, "/api/v1/ping", nil))

	if rr.Code != http.StatusUnauthorized {
		t.Fatalf("status = %d, want %d", rr.Code, http.StatusUnauthorized)
	}
	if rr.Header().Get("WWW-Authenticate") != "Bearer" {
		t.Fatalf("WWW-Authenticate = %q", rr.Header().Get("WWW-Authenticate"))
	}
}
