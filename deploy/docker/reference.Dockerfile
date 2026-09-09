# Builds implementations/reference/cmd/server. See deploy/README.md — local development only,
# not a production deployment artifact.
#
# Build from the repo root:
#   docker build -f deploy/docker/reference.Dockerfile -t bongopay-reference .

FROM golang:1.27-alpine AS build
WORKDIR /src

# No go.sum: the reference implementation has zero external dependencies by design (ADR 0002),
# so there's nothing for `go mod download` to cache separately from the source copy below.
COPY implementations/reference/go.mod ./implementations/reference/go.mod
COPY implementations/reference/ ./implementations/reference/

WORKDIR /src/implementations/reference
RUN CGO_ENABLED=0 go build -o /out/bongopay-server ./cmd/server

FROM alpine:3
RUN adduser -D -H bongopay
USER bongopay
COPY --from=build /out/bongopay-server /usr/local/bin/bongopay-server

EXPOSE 8080
ENTRYPOINT ["/usr/local/bin/bongopay-server"]
