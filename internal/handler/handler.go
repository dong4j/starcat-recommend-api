// Package handler 提供 HTTP handler 公共工具。
package handler

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/dong4j/starcat-recommend-api/internal/model"
)

func writeJSON[T any](w http.ResponseWriter, data T) {
	writeJSONWithMeta(w, data, nil)
}

func writeJSONWithMeta[T any](w http.ResponseWriter, data T, meta *model.Meta) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)

	env := model.Envelope[T]{
		SchemaVersion: 1,
		Data:          data,
		Meta:          meta,
	}
	if err := json.NewEncoder(w).Encode(env); err != nil {
		log.Printf("[handler] failed to encode envelope: %v", err)
	}
}

func writeError(w http.ResponseWriter, status int, code, message string, details interface{}) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)

	env := model.ErrorEnvelope{
		SchemaVersion: 1,
		Error: model.ErrorResponse{
			Code:    code,
			Message: message,
			Details: details,
		},
	}
	if err := json.NewEncoder(w).Encode(env); err != nil {
		log.Printf("[handler] failed to encode error envelope: %v", err)
	}
}
