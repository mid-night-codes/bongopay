# deploy/compose/ — Docker Compose Files

Docker Compose files to run BongoPay components together locally (e.g. the reference
implementation together with the simulator), once the images built from
[deploy/docker/](../docker/) exist. See [deploy/README.md](../README.md) for the rules governing
this directory.

## Status

**`reference.yml`** — a single service running
[deploy/docker/reference.Dockerfile](../docker/reference.Dockerfile), exposing `:8080`. Only one
component exists to compose so far ("compose" here just means "build and run one image
consistently"; this becomes more useful once there's more than one service to wire together).

```bash
docker compose -f deploy/compose/reference.yml up --build
curl -X POST localhost:8080/payments \
  -d '{"provider":{"id":"SIMULATOR"},"amount":{"value":5000,"currency":{"code":"TZS"}},"customerReference":{},"idempotencyKey":"demo-1"}'
```

The callback-signing secret is set to an obviously-fake, fixed value via `command:` (see
[AGENTS.md §10](../../AGENTS.md#10-security-restrictions)) purely so a signature computed
against it stays valid across container restarts during local testing — never a real
credential, and never meaningful outside this local-dev context.
