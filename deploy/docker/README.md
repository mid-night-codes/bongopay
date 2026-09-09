# deploy/docker/ — Dockerfiles

Dockerfile(s) for running BongoPay components locally, once there are components to run. See
[deploy/README.md](../README.md) for the rules governing this directory (local development only,
no real credentials, reflects `implementations/`/`adapters/` rather than leading them).

## Status

**`reference.Dockerfile`** — a multi-stage build (`golang:1.27-alpine` builder, `alpine:3`
runtime, non-root user) for
[implementations/reference/cmd/server](../../implementations/reference/README.md). No
`go.sum` to copy separately — the reference implementation has zero external dependencies by
design (see [ADR 0002](../../adr/0002-reference-implementation-language-go.md)). Build from the
repo root:

```bash
docker build -f deploy/docker/reference.Dockerfile -t bongopay-reference .
docker run --rm -p 8080:8080 bongopay-reference
```

Prefer [deploy/compose/reference.yml](../compose/reference.yml) for everyday local use — it
wraps this same build.
