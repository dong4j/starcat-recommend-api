#!/usr/bin/env bash
set -euo pipefail

VERSION="${1:-}"
if [[ -z "${VERSION}" ]]; then
  echo "usage: ./scripts/deploy.sh vX.Y.Z"
  exit 1
fi

go test ./...
go test -race ./...
go vet ./...
go build ./...

git tag "${VERSION}"
git push origin "${VERSION}"
