package provider

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestSimRepoProviderMapsResponse(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Fatalf("method = %s", r.Method)
		}
		if got := r.Header.Get("api-key"); got != "sk-simrepo" {
			t.Fatalf("api-key = %q", got)
		}

		var req simRepoRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			t.Fatalf("decode request: %v", err)
		}
		if req.Limit != 1 || req.Offset != 2 || len(req.Positive) != 1 || req.Positive[0] != 41881900 {
			t.Fatalf("unexpected request: %+v", req)
		}

		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{
			"result": [
				{
					"id": 123,
					"score": 0.87,
					"payload": {
						"full_name": "owner/repo",
						"description": "desc",
						"language": "Go",
						"stargazers_count": 1200,
						"forks_count": 80,
						"archived": false
					}
				}
			],
			"status": "ok",
			"time": 0.001
		}`))
	}))
	defer server.Close()

	provider := NewSimRepoProvider(server.URL, "sk-simrepo", server.Client())
	result, err := provider.Recommend(t.Context(), Query{RepoID: 41881900, Limit: 1, Offset: 2})
	if err != nil {
		t.Fatalf("Recommend: %v", err)
	}

	if result.Response.Source != "simrepo" || result.Response.RepoID != 41881900 {
		t.Fatalf("unexpected response header: %+v", result.Response)
	}
	if len(result.Response.Items) != 1 {
		t.Fatalf("items len = %d", len(result.Response.Items))
	}
	item := result.Response.Items[0]
	if item.RepoID != 123 || item.FullName != "owner/repo" || item.Stars != 1200 || item.Forks != 80 {
		t.Fatalf("unexpected item: %+v", item)
	}
	if !result.Response.HasMore || result.Response.NextOffset == nil || *result.Response.NextOffset != 3 {
		t.Fatalf("unexpected paging: %+v", result.Response)
	}
}
