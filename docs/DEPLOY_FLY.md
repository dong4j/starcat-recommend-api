# Deploy to Fly.io

```bash
fly launch --no-deploy
fly secrets set API_KEYS="sk-starcat-..."
fly secrets set SIMREPO_API_KEY="..."
fly deploy
```

Health check:

```bash
curl https://starcat-recommend-api.fly.dev/healthz
```
