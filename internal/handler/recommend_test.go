package handler

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/dong4j/starcat-recommend-api/internal/model"
	"github.com/dong4j/starcat-recommend-api/internal/provider"
)

type mockRecommendationProvider struct {
	result provider.Result
	err    error
	query  provider.Query
}

func (m *mockRecommendationProvider) Recommend(ctx context.Context, query provider.Query) (provider.Result, error) {
	m.query = query
	return m.result, m.err
}

func TestRecommendHandlerSuccess(t *testing.T) {
	mock := &mockRecommendationProvider{
		result: provider.Result{
			Response: model.RecommendationResponse{
				Source: "simrepo",
				RepoID: 41881900,
				Items: []model.RecommendationItem{{
					RepoID:   1,
					FullName: "owner/repo",
					Score:    0.9,
					Source:   "simrepo",
				}},
			},
			CacheStatus: "fresh",
		},
	}
	req := httptest.NewRequest(http.MethodGet, "/api/v1/repos/41881900/recommendations?limit=10&offset=20", nil)
	req.SetPathValue("repo_id", "41881900")
	rr := httptest.NewRecorder()

	NewRecommendHandler(mock).HandleRecommendations(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d", rr.Code, http.StatusOK)
	}
	if mock.query.RepoID != 41881900 || mock.query.Limit != 10 || mock.query.Offset != 20 {
		t.Fatalf("query = %+v", mock.query)
	}

	var env model.Envelope[model.RecommendationResponse]
	if err := json.NewDecoder(rr.Body).Decode(&env); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if len(env.Data.Items) != 1 || env.Data.Items[0].FullName != "owner/repo" {
		t.Fatalf("unexpected response: %+v", env.Data)
	}
	if env.Meta == nil || env.Meta.CacheStatus != "fresh" {
		t.Fatalf("unexpected meta: %+v", env.Meta)
	}
}

func TestRecommendHandlerBadRepoID(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/api/v1/repos/abc/recommendations", nil)
	req.SetPathValue("repo_id", "abc")
	rr := httptest.NewRecorder()

	NewRecommendHandler(&mockRecommendationProvider{}).HandleRecommendations(rr, req)

	if rr.Code != http.StatusBadRequest {
		t.Fatalf("status = %d, want %d", rr.Code, http.StatusBadRequest)
	}
}

func TestRecommendHandlerUpstreamError(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/api/v1/repos/1/recommendations", nil)
	req.SetPathValue("repo_id", "1")
	rr := httptest.NewRecorder()

	NewRecommendHandler(&mockRecommendationProvider{err: errors.New("boom")}).HandleRecommendations(rr, req)

	if rr.Code != http.StatusBadGateway {
		t.Fatalf("status = %d, want %d", rr.Code, http.StatusBadGateway)
	}
}
