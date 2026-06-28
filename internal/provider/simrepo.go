package provider

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/dong4j/starcat-recommend-api/internal/model"
)

const simRepoSource = "simrepo"

// SimRepoProvider 调用 SimRepo 暴露的 Qdrant recommend endpoint。
//
// 这里故意只输出 Starcat 自己的 RecommendationItem, 不把 Qdrant 原始字段透传出去。
type SimRepoProvider struct {
	endpoint string
	apiKey   string
	client   *http.Client
}

func NewSimRepoProvider(endpoint, apiKey string, client *http.Client) *SimRepoProvider {
	if client == nil {
		client = &http.Client{Timeout: 20 * time.Second}
	}
	return &SimRepoProvider{
		endpoint: endpoint,
		apiKey:   apiKey,
		client:   client,
	}
}

func (p *SimRepoProvider) Recommend(ctx context.Context, query Query) (Result, error) {
	if query.RepoID <= 0 {
		return Result{}, errors.New("repo_id must be positive")
	}

	body := simRepoRequest{
		Limit:       query.Limit,
		Positive:    []int64{query.RepoID},
		Filter:      simRepoFilter{Must: []any{}},
		Offset:      query.Offset,
		WithPayload: true,
		WithVector:  false,
	}
	payload, err := json.Marshal(body)
	if err != nil {
		return Result{}, err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, p.endpoint, bytes.NewReader(payload))
	if err != nil {
		return Result{}, err
	}
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("User-Agent", "starcat-recommend-api/1.0")
	req.Header.Set("api-key", p.apiKey)

	resp, err := p.client.Do(req)
	if err != nil {
		return Result{}, fmt.Errorf("simrepo request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode > 299 {
		return Result{}, fmt.Errorf("simrepo returned HTTP %d", resp.StatusCode)
	}

	var decoded simRepoResponse
	decoder := json.NewDecoder(resp.Body)
	decoder.UseNumber()
	if err := decoder.Decode(&decoded); err != nil {
		return Result{}, fmt.Errorf("simrepo response decode failed: %w", err)
	}

	items := make([]model.RecommendationItem, 0, len(decoded.Result))
	for _, point := range decoded.Result {
		item := mapPoint(point)
		if item.RepoID <= 0 || item.FullName == "" {
			continue
		}
		items = append(items, item)
	}

	var nextOffset *int
	hasMore := len(items) == query.Limit && query.Limit > 0
	if hasMore {
		next := query.Offset + query.Limit
		nextOffset = &next
	}

	return Result{
		Response: model.RecommendationResponse{
			Source:     simRepoSource,
			Fallback:   false,
			RepoID:     query.RepoID,
			Items:      items,
			HasMore:    hasMore,
			NextOffset: nextOffset,
		},
		CacheStatus: "fresh",
	}, nil
}

type simRepoRequest struct {
	Limit       int           `json:"limit"`
	Positive    []int64       `json:"positive"`
	Filter      simRepoFilter `json:"filter"`
	Offset      int           `json:"offset"`
	WithPayload bool          `json:"with_payload"`
	WithVector  bool          `json:"with_vector"`
}

type simRepoFilter struct {
	Must []any `json:"must"`
}

type simRepoResponse struct {
	Result []simRepoPoint `json:"result"`
	Status string         `json:"status"`
	Time   float64        `json:"time"`
}

type simRepoPoint struct {
	ID      any            `json:"id"`
	Score   float64        `json:"score"`
	Payload map[string]any `json:"payload"`
}

func mapPoint(point simRepoPoint) model.RecommendationItem {
	repoID := int64Value(point.Payload, "id")
	if repoID <= 0 {
		repoID = anyInt64(point.ID)
	}

	return model.RecommendationItem{
		RepoID:      repoID,
		FullName:    stringValue(point.Payload, "full_name"),
		Description: stringValue(point.Payload, "description"),
		Language:    stringValue(point.Payload, "language"),
		Stars:       intValue(point.Payload, "stargazers_count", "stars"),
		Forks:       intValue(point.Payload, "forks_count", "forks"),
		Archived:    boolValue(point.Payload, "archived"),
		Score:       point.Score,
		Source:      simRepoSource,
		Reasons:     []string{"被相似 GitHub 用户共同 star"},
	}
}

func stringValue(payload map[string]any, key string) string {
	value, ok := payload[key]
	if !ok || value == nil {
		return ""
	}
	switch typed := value.(type) {
	case string:
		return typed
	default:
		return strings.TrimSpace(fmt.Sprint(typed))
	}
}

func intValue(payload map[string]any, keys ...string) int {
	for _, key := range keys {
		if value := int64Value(payload, key); value > 0 {
			return int(value)
		}
	}
	return 0
}

func int64Value(payload map[string]any, key string) int64 {
	value, ok := payload[key]
	if !ok {
		return 0
	}
	return anyInt64(value)
}

func anyInt64(value any) int64 {
	switch typed := value.(type) {
	case json.Number:
		parsed, _ := typed.Int64()
		return parsed
	case float64:
		return int64(typed)
	case int64:
		return typed
	case int:
		return int64(typed)
	case string:
		parsed, _ := strconv.ParseInt(typed, 10, 64)
		return parsed
	default:
		return 0
	}
}

func boolValue(payload map[string]any, key string) bool {
	value, ok := payload[key]
	if !ok {
		return false
	}
	switch typed := value.(type) {
	case bool:
		return typed
	case string:
		return typed == "true"
	default:
		return false
	}
}
