# starcat-recommend-api

<!-- starcat-promo:start -->
<div align="center">
<a href="https://starcat.ink"><img src="https://raw.githubusercontent.com/starcat-app/starcat-pro/main/banner.webp" width="100%" alt="Starcat" /></a>

<p><strong>Self-hostable support API for Starcat similar-repository recommendations.</strong></p>
<p>Starcat is a native macOS app that turns GitHub Stars into a searchable, organized and AI-assisted knowledge base. It supports README rendering, tags, private notes, release tracking, repository health signals, AI summaries, semantic search, browser plugin workflows and self-hostable support APIs.</p>

<a href="https://github.com/starcat-app/homebrew-starcat"><img src="https://img.shields.io/badge/Install%20with-Homebrew-FBBF24?style=for-the-badge&logo=homebrew&logoColor=white" width="220" alt="Install with Homebrew"/></a>
<br/>
<sub><a href="./README-ZH.md">中文说明</a></sub>
</div>

<div align="center">
<a href="https://starcat.ink"><img src="https://img.shields.io/badge/website-starcat.ink-38BDF8?style=flat&color=blue" alt="website"/></a>
<a href="https://github.com/starcat-app/starcat-pro"><img src="https://img.shields.io/badge/support-starcat--pro-lightgrey.svg?style=flat&color=blue" alt="support"/></a>
<a href="https://github.com/starcat-app/homebrew-starcat"><img src="https://img.shields.io/badge/install-homebrew-lightgrey.svg?style=flat&color=blue" alt="homebrew"/></a>
<a href="https://github.com/starcat-app/starcat-localization"><img src="https://img.shields.io/badge/localization-open-lightgrey.svg?style=flat&color=blue" alt="localization"/></a>
</div>

<div align="center">
<img width="900" src="https://raw.githubusercontent.com/starcat-app/starcat-pro/main/main.webp" alt="Starcat main window"/>
</div>

**Preferred install method:**

```bash
brew tap starcat-app/starcat
brew trust starcat-app/starcat
brew install --cask starcat
```

**Useful links:**

- Home: https://starcat.ink
- Download: https://starcat.ink/downloads/Starcat-1.1.0-arm64.dmg
- Public support and release notes: https://github.com/starcat-app/starcat-pro
- Homebrew tap: https://github.com/starcat-app/homebrew-starcat
- Browser plugins: [Chrome](https://github.com/starcat-app/starcat-chrome-plugin) / [Safari](https://github.com/starcat-app/starcat-safari-plugin)
- Localization: https://github.com/starcat-app/starcat-localization

**Starcat ecosystem:**

- [starcat-sharing-api](https://github.com/dong4j/starcat-sharing-api)
- [starcat-trending-api](https://github.com/dong4j/starcat-trending-api)
- [starcat-weekly-api](https://github.com/dong4j/starcat-weekly-api)
- [starcat-wiki-api](https://github.com/dong4j/starcat-wiki-api)
- [starcat-recommend-api](https://github.com/dong4j/starcat-recommend-api)
- [starcat-discovery-api](https://github.com/dong4j/starcat-discovery-api)
- [starcat-license-api](https://github.com/dong4j/starcat-license-api)

> Starcat provides hosted defaults for normal users. This API is open source so advanced users can inspect it, run it locally, or deploy their own instance.
<!-- starcat-promo:end -->

Backend service for Starcat similar-repository recommendations.

The initial version proxies SimRepo's unofficial Qdrant Recommend API through the server, providing the Starcat macOS client with a stable `/api/v1/repos/{repo_id}/recommendations` contract. The client neither connects directly to SimRepo nor stores a SimRepo key. The provider can later be replaced within this service by Starcat's own recommendation provider.

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
- `CACHE_TTL_EMPTY_SECONDS`: defaults to 1 hour.
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

`CachedProvider` keeps at most 10,000 `repoID:limit:offset` entries. Expired entries are removed on read; when capacity is reached, the entry with the earliest expiry is evicted.

Future providers should keep the response DTO stable:

```text
ContentEmbeddingProvider
StarcatBehaviorProvider
HybridProvider
```
