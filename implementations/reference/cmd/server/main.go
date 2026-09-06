// Command server runs the BongoPay reference implementation's HTTP server, implementing
// contracts/openapi/bongopay.yaml against an in-memory Store and the SIMULATOR provider — see
// implementations/reference/README.md. It is a reference/demo entrypoint, not a deployable
// production server: persistence is in-memory only (ARCHITECTURE.md §14, "Persistence
// architecture" is a TODO(ADR)), and there is no authentication (see bongopay.yaml's own
// TODO(spec) on that).
package main

import (
	"crypto/rand"
	"encoding/hex"
	"flag"
	"log"
	"net/http"

	"github.com/mid-night-codes/bongopay/implementations/reference/internal/httpapi"
	"github.com/mid-night-codes/bongopay/implementations/reference/internal/payment"
	"github.com/mid-night-codes/bongopay/implementations/reference/internal/simulator"
)

func main() {
	addr := flag.String("addr", ":8080", "address to listen on")
	callbackSecret := flag.String("callback-secret", "", "HMAC secret for simulator callback signatures (random if empty)")
	flag.Parse()

	secret := []byte(*callbackSecret)
	if len(secret) == 0 {
		var b [32]byte
		if _, err := rand.Read(b[:]); err != nil {
			log.Fatalf("generating callback secret: %v", err)
		}
		secret = []byte(hex.EncodeToString(b[:]))
	}

	store := payment.NewInMemoryStore()
	svc := payment.NewService(store)
	sim := simulator.New(svc, simulator.WithSecret(secret))
	server := httpapi.NewServer(svc, sim)

	log.Printf("bongopay reference server listening on %s", *addr)
	log.Printf("simulator callback-signing secret (POST /simulator/callbacks, header %s): %s",
		httpapi.SignatureHeader, secret)
	if err := http.ListenAndServe(*addr, server.Handler()); err != nil {
		log.Fatal(err)
	}
}
