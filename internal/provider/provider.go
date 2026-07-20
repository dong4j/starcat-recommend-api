// Package provider 封装推荐数据源。
//
// handler 只依赖本包的 Query/Result; 后续从 SimRepo 切到 Starcat 自研推荐时,
// 保持 Result 语义不变即可避免客户端契约变化。
package provider

import (
	"context"

	"github.com/starcat-app/starcat-recommend-api/internal/model"
)

type Query struct {
	RepoID int64
	Limit  int
	Offset int
}

type Result struct {
	Response    model.RecommendationResponse
	CacheStatus string
}

type Provider interface {
	Recommend(ctx context.Context, query Query) (Result, error)
}
