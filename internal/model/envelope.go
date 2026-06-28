// Package model 定义 Envelope 统一响应结构。
//
// 与 Starcat 现有自建 API 保持同款 envelope, 让客户端可复用统一解码路径。
package model

// Envelope 是 /api/v1/* 200 响应的顶层包装。
type Envelope[T any] struct {
	SchemaVersion int   `json:"schema_version"`
	Data          T     `json:"data"`
	Meta          *Meta `json:"meta,omitempty"`
}

// Meta 可选的分页/缓存/来源元数据。
type Meta struct {
	Page        int    `json:"page,omitempty"`
	PageSize    int    `json:"page_size,omitempty"`
	Total       int    `json:"total,omitempty"`
	NextPage    *int   `json:"next_page,omitempty"`
	Source      string `json:"source,omitempty"`
	CacheStatus string `json:"cache_status,omitempty"`
	FetchedAt   string `json:"fetched_at,omitempty"`
}

// ErrorResponse 统一错误响应体。
type ErrorResponse struct {
	Code    string      `json:"code"`
	Message string      `json:"message"`
	Details interface{} `json:"details,omitempty"`
}

// ErrorEnvelope 所有非 2xx 响应的顶层包装。
type ErrorEnvelope struct {
	SchemaVersion int           `json:"schema_version"`
	Error         ErrorResponse `json:"error"`
}
