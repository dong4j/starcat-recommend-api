// Package middleware 提供 Bearer Token 鉴权中间件。
package middleware

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"

	"github.com/starcat-app/starcat-recommend-api/internal/model"
)

// BearerAuth 持有 API Key 白名单, 验证 Bearer Token。
type BearerAuth struct {
	allowedKeys map[string]bool
}

// NewBearerAuth 创建 Bearer 鉴权中间件。
func NewBearerAuth(keys []string) *BearerAuth {
	m := make(map[string]bool, len(keys))
	for _, k := range keys {
		k = strings.TrimSpace(k)
		if k != "" {
			m[k] = true
		}
	}
	log.Printf("[auth] %d keys loaded", len(m))
	return &BearerAuth{allowedKeys: m}
}

// Wrap 返回一个 http.Handler, 在执行业务 handler 前验证 Bearer Token。
func (a *BearerAuth) Wrap(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")

		if authHeader == "" {
			writeAuthError(w, "missing Authorization header")
			return
		}

		if !strings.HasPrefix(authHeader, "Bearer ") {
			writeAuthError(w, "expected 'Bearer <token>' format")
			return
		}

		token := strings.TrimSpace(strings.TrimPrefix(authHeader, "Bearer "))
		if !a.allowedKeys[token] {
			log.Printf("[auth] rejected key %s", maskKey(token))
			writeAuthError(w, "invalid API key")
			return
		}

		next.ServeHTTP(w, r)
	})
}

func writeAuthError(w http.ResponseWriter, msg string) {
	w.Header().Set("WWW-Authenticate", "Bearer")
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusUnauthorized)
	_ = json.NewEncoder(w).Encode(model.ErrorEnvelope{
		SchemaVersion: 1,
		Error: model.ErrorResponse{
			Code:    "UNAUTHORIZED",
			Message: msg,
		},
	})
}

func maskKey(key string) string {
	if len(key) < 16 {
		return "****"
	}
	return key[:7] + "****" + key[len(key)-4:]
}
