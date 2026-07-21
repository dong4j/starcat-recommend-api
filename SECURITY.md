# Security Policy

## Reporting a vulnerability

Report authentication bypasses, SSRF, unsafe upstream responses, cache poisoning, or secret exposure through [GitHub Security Advisories](https://github.com/starcat-app/starcat-recommend-api/security/advisories/new). Do not publish API keys, upstream credentials, production logs, or user repository data in an issue.

Security fixes target the current default branch and latest deployed version. Runtime secrets must be injected through environment variables or Fly.io secrets and must never be committed.
