package handler

import (
	"context"
	"net/http"
	"strconv"

	"github.com/starcat-app/starcat-recommend-api/internal/model"
	"github.com/starcat-app/starcat-recommend-api/internal/provider"
)

const (
	defaultLimit = 10
	maxLimit     = 30
)

// RecommendationProvider 是 handler 依赖的最小 provider 接口。
type RecommendationProvider interface {
	Recommend(ctx context.Context, query provider.Query) (provider.Result, error)
}

// RecommendHandler 处理相似仓库推荐接口。
type RecommendHandler struct {
	provider RecommendationProvider
}

func NewRecommendHandler(provider RecommendationProvider) *RecommendHandler {
	return &RecommendHandler{provider: provider}
}

func (h *RecommendHandler) HandleRecommendations(w http.ResponseWriter, r *http.Request) {
	repoID, ok := parsePositiveInt64(r.PathValue("repo_id"))
	if !ok {
		writeError(w, http.StatusBadRequest, "BAD_REQUEST", "repo_id must be positive", nil)
		return
	}

	limit, ok := parseBoundedInt(r.URL.Query().Get("limit"), defaultLimit, 1, maxLimit)
	if !ok {
		writeError(w, http.StatusBadRequest, "BAD_REQUEST", "limit must be between 1 and 30", nil)
		return
	}

	offset, ok := parseBoundedInt(r.URL.Query().Get("offset"), 0, 0, 10000)
	if !ok {
		writeError(w, http.StatusBadRequest, "BAD_REQUEST", "offset must be non-negative", nil)
		return
	}

	result, err := h.provider.Recommend(r.Context(), provider.Query{
		RepoID: repoID,
		Limit:  limit,
		Offset: offset,
	})
	if err != nil {
		writeError(w, http.StatusBadGateway, "UPSTREAM_UNAVAILABLE", err.Error(), nil)
		return
	}

	meta := &model.Meta{
		PageSize:    limit,
		Total:       len(result.Response.Items),
		Source:      result.Response.Source,
		CacheStatus: result.CacheStatus,
	}
	writeJSONWithMeta(w, result.Response, meta)
}

func parsePositiveInt64(raw string) (int64, bool) {
	value, err := strconv.ParseInt(raw, 10, 64)
	if err != nil || value <= 0 {
		return 0, false
	}
	return value, true
}

func parseBoundedInt(raw string, fallback, min, max int) (int, bool) {
	if raw == "" {
		return fallback, true
	}
	value, err := strconv.Atoi(raw)
	if err != nil || value < min || value > max {
		return 0, false
	}
	return value, true
}
