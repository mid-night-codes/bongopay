// Package testcontainer boots implementations/reference as a Docker container for tests, using
// testcontainers-go — see deploy/docker/reference.Dockerfile for the image this builds.
//
// This is a separate, nested Go module (its own go.mod) from sdks/go on purpose: it's the first
// external dependency in this project's history (everything else is stdlib-only by design, per
// ADR 0002), and someone importing sdks/go to make API calls has no reason to also pull in
// Docker-orchestration machinery. Only importing this package takes on that dependency.
package testcontainer

import (
	"context"
	"fmt"

	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

// referencePort is the port cmd/server listens on inside the container — see
// deploy/docker/reference.Dockerfile's EXPOSE.
const referencePort = "8080/tcp"

// ReferenceContainer is a running reference implementation container, ready to accept requests.
type ReferenceContainer struct {
	testcontainers.Container

	// BaseURL is the container's address from the host's point of view — e.g.
	// "http://localhost:54321" — suitable for an sdks/go bongopay.New(...) call.
	BaseURL string
}

// RunReference builds deploy/docker/reference.Dockerfile and starts it, waiting for cmd/server's
// own startup log line before returning. repoRoot is the BongoPay repository root — the image's
// build context (see deploy/docker/README.md); callers typically resolve this relative to their
// own test file. Callers must call Terminate (or defer it) to stop the container.
func RunReference(ctx context.Context, repoRoot string) (*ReferenceContainer, error) {
	req := testcontainers.ContainerRequest{
		FromDockerfile: testcontainers.FromDockerfile{
			Context:    repoRoot,
			Dockerfile: "deploy/docker/reference.Dockerfile",
		},
		ExposedPorts: []string{referencePort},
		// wait.ForLog alone isn't enough here: cmd/server logs its startup line just
		// *before* calling http.ListenAndServe, and even once it's genuinely listening,
		// Docker's own port-mapping/proxy can lag slightly behind — verified by this
		// package's own test, which flaked with a connection EOF using ForLog alone.
		// ForListeningPort checks the actual host-mapped port is reachable, which is the
		// real thing callers need (they'll connect the same way), not just that the
		// process printed something.
		WaitingFor: wait.ForListeningPort(referencePort),
	}

	container, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	if err != nil {
		return nil, fmt.Errorf("testcontainer: starting reference container: %w", err)
	}

	host, err := container.Host(ctx)
	if err != nil {
		return nil, fmt.Errorf("testcontainer: resolving container host: %w", err)
	}
	port, err := container.MappedPort(ctx, referencePort)
	if err != nil {
		return nil, fmt.Errorf("testcontainer: resolving mapped port: %w", err)
	}

	return &ReferenceContainer{
		Container: container,
		BaseURL:   fmt.Sprintf("http://%s:%s", host, port.Port()),
	}, nil
}
