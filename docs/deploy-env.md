# Deploy Environment

Required secrets:

```bash
fly secrets set API_KEYS="sk-starcat-..."
fly secrets set SIMREPO_API_KEY="..."
```

Optional settings:

```bash
fly secrets set SIMREPO_ENDPOINT="https://simrepo.dera.page/collections/repos/points/recommend"
fly secrets set CACHE_TTL_SUCCESS_SECONDS="604800"
fly secrets set CACHE_TTL_EMPTY_SECONDS="3600"
fly secrets set CACHE_TTL_ERROR_SECONDS="600"
```
