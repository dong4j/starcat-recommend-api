package provider

import (
	"context"
	"testing"
	"time"

	"github.com/dong4j/starcat-recommend-api/internal/model"
)

type countingProvider struct {
	calls int
}

func (p *countingProvider) Recommend(ctx context.Context, query Query) (Result, error) {
	p.calls++
	return Result{
		Response: model.RecommendationResponse{
			Source: "simrepo",
			RepoID: query.RepoID,
			Items:  []model.RecommendationItem{{RepoID: 1, FullName: "owner/repo"}},
		},
		CacheStatus: "fresh",
	}, nil
}

func TestCachedProviderCachesSuccess(t *testing.T) {
	base := &countingProvider{}
	cached := NewCachedProvider(base, time.Minute, time.Minute, time.Minute)
	query := Query{RepoID: 1, Limit: 10, Offset: 0}

	first, err := cached.Recommend(t.Context(), query)
	if err != nil {
		t.Fatalf("first: %v", err)
	}
	second, err := cached.Recommend(t.Context(), query)
	if err != nil {
		t.Fatalf("second: %v", err)
	}

	if base.calls != 1 {
		t.Fatalf("base calls = %d, want 1", base.calls)
	}
	if first.CacheStatus != "fresh" || second.CacheStatus != "hit" {
		t.Fatalf("cache statuses = %q/%q", first.CacheStatus, second.CacheStatus)
	}
}
