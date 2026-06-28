# Custom Domain

The Starcat client stores only the service base URL, for example:

```text
https://recommend.example.com
```

Do not include `/api` in the base URL. API paths are absolute and are appended by the client:

```text
/api/v1/ping
/api/v1/repos/{repo_id}/recommendations
```
