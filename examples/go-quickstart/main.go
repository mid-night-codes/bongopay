// Command go-quickstart demonstrates the flow that's actually supported end-to-end via the
// public HTTP API today: initiate a payment (success and failure scenarios) and read it back.
// See examples/README.md for what this directory is for and why a "payment + webhook" example
// — the shape that README suggests as a model example — isn't built yet: Simulator.Initiate
// resolves success/failure synchronously in one call rather than through a separate callback
// step, which is an open design question recorded in ROADMAP.md, not something this example
// should paper over with an invented demonstration.
package main

import (
	"context"
	"flag"
	"fmt"
	"log"

	bongopay "github.com/mid-night-codes/bongopay/sdks/go"
	"github.com/mid-night-codes/bongopay/sdks/go/generated"
)

func main() {
	server := flag.String("server", "http://localhost:8080", "reference server base URL")
	flag.Parse()

	client := bongopay.New(*server)
	ctx := context.Background()

	fmt.Println("== Initiating a payment against the SIMULATOR (success scenario) ==")
	success, err := client.InitiatePayment(ctx, request("quickstart-success", "success"))
	if err != nil {
		log.Fatalf("initiate (success) failed: %v\n\nIs the reference server running? See examples/go-quickstart/README.md.", err)
	}
	printPayment(success.Payment)

	fmt.Println("\n== Initiating a payment against the SIMULATOR (failure scenario) ==")
	failure, err := client.InitiatePayment(ctx, request("quickstart-failure", "failure"))
	if err != nil {
		log.Fatalf("initiate (failure) failed: %v", err)
	}
	printPayment(failure.Payment)

	fmt.Println("\n== Reading the first payment back by ID ==")
	got, err := client.GetPayment(ctx, success.Payment.Id)
	if err != nil {
		log.Fatalf("get failed: %v", err)
	}
	printPayment(*got)

	fmt.Println("\nSame idempotency key returns the same payment, not a new one:")
	replay, err := client.InitiatePayment(ctx, request("quickstart-success", "success"))
	if err != nil {
		log.Fatalf("replay initiate failed: %v", err)
	}
	if replay.Payment.Id != success.Payment.Id {
		log.Fatalf("expected the same payment ID on replay, got a different one: %s vs %s", replay.Payment.Id, success.Payment.Id)
	}
	fmt.Printf("Replayed idempotency key %q -> same payment ID %s\n", "quickstart-success", replay.Payment.Id)
}

func request(idempotencyKey, scenario string) generated.PaymentRequest {
	return generated.PaymentRequest{
		Provider:          generated.Provider{Id: "SIMULATOR"},
		Amount:            generated.Money{Value: 5000, Currency: generated.Currency{Code: "TZS"}},
		CustomerReference: generated.CustomerReference{},
		IdempotencyKey:    idempotencyKey,
		ProviderOptions: &generated.ProviderOptions{
			"simulator": map[string]interface{}{"scenario": scenario},
		},
	}
}

func printPayment(p generated.Payment) {
	fmt.Printf("  id:              %s\n", p.Id)
	fmt.Printf("  status:          %s\n", p.Status)
	fmt.Printf("  amount:          %d %s\n", p.Amount.Value, p.Amount.Currency.Code)
	fmt.Printf("  idempotencyKey:  %s\n", p.IdempotencyKey)
}
