// Command server runs the BongoPay reference implementation's HTTP server, implementing
// contracts/openapi/bongopay.yaml against an in-memory Store and the SIMULATOR provider — see
// implementations/reference/README.md. It is a reference/demo entrypoint, not a deployable
// production server: persistence is in-memory only (ARCHITECTURE.md §14, "Persistence
// architecture" is a TODO(ADR)), and there is no authentication (see bongopay.yaml's own
// TODO(spec) on that).
package main

import (
	"flag"
	"log"
	"net/http"

	"github.com/mid-night-codes/bongopay/implementations/reference/internal/httpapi"
	"github.com/mid-night-codes/bongopay/implementations/reference/internal/payment"
	"github.com/mid-night-codes/bongopay/implementations/reference/internal/simulator"
)

func main() {
	addr := flag.String("addr", ":8080", "address to listen on")
	flag.Parse()

	store := payment.NewInMemoryStore()
	svc := payment.NewService(store)
	sim := simulator.New(svc)
	server := httpapi.NewServer(svc, sim)

	log.Printf("bongopay reference server listening on %s", *addr)
	if err := http.ListenAndServe(*addr, server.Handler()); err != nil {
		log.Fatal(err)
	}
}
