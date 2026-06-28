# starcat-recommend-api

Starcat 相似仓库推荐后端。

第一版通过服务端中转 SimRepo 的非官方 Qdrant Recommend API, 为 Starcat macOS 客户端提供稳定的 `/api/v1/repos/{repo_id}/recommendations` 契约。客户端不直连 SimRepo, 不持有 SimRepo key, 后续可以在本服务内替换为 Starcat 自研推荐 Provider。

## Endpoints

| Method | Path | Auth | Description |
|---|---|---:|---|
| `GET` | `/healthz` | No | Process health check |
| `GET` | `/api/v1/ping` | Yes | Starcat client connectivity probe |
| `GET` | `/api/v1/repos/{repo_id}/recommendations?limit=10&offset=0` | Yes | Similar repository recommendations |

## Environment

```bash
cp .env.example .env
```

Required:

- `API_KEYS`: comma-separated Bearer tokens accepted by Starcat clients.
- `SIMREPO_API_KEY`: SimRepo Qdrant read-only key. Keep it server-side only.

Optional:

- `PORT`: defaults to `5005`.
- `SIMREPO_ENDPOINT`: defaults to SimRepo's Qdrant recommend endpoint.
- `CACHE_TTL_SUCCESS_SECONDS`: defaults to 7 days.
- `CACHE_TTL_EMPTY_SECONDS`: defaults to 6 hours.
- `CACHE_TTL_ERROR_SECONDS`: defaults to 10 minutes.

## Local Development

```bash
go mod tidy
go run ./cmd/server
```

Smoke test:

```bash
curl http://127.0.0.1:5005/healthz
curl -H "Authorization: Bearer $API_KEY" http://127.0.0.1:5005/api/v1/ping
curl -H "Authorization: Bearer $API_KEY" \
  "http://127.0.0.1:5005/api/v1/repos/41881900/recommendations?limit=10&offset=0"
```

## Quality Gates

```bash
go test ./...
go test -race ./...
go vet ./...
go build ./...
```

## Provider Boundary

The current provider chain is:

```text
RecommendHandler -> CachedProvider -> SimRepoProvider -> SimRepo Qdrant API
```

Future providers should keep the response DTO stable:

```text
ContentEmbeddingProvider
StarcatBehaviorProvider
HybridProvider
```
