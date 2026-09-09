// Command cli is a small command-line client for the reference implementation's HTTP server
// (cmd/server) — see implementations/reference/README.md. It exists so `initiate`/`get` can be
// exercised without hand-writing curl/jq invocations; it is not a general-purpose BongoPay
// product CLI (that's a separate, undecided, later scope — see ROADMAP.md Phase 2).
package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/mid-night-codes/bongopay/implementations/reference/internal/payment"
)

func main() {
	if len(os.Args) < 2 {
		usage()
		os.Exit(2)
	}

	switch os.Args[1] {
	case "initiate":
		runInitiate(os.Args[2:])
	case "get":
		runGet(os.Args[2:])
	case "-h", "--help", "help":
		usage()
	default:
		fmt.Fprintf(os.Stderr, "unknown command %q\n\n", os.Args[1])
		usage()
		os.Exit(2)
	}
}

func usage() {
	fmt.Fprintln(os.Stderr, `bongopay-cli — talks to the reference implementation's HTTP server (cmd/server)

Usage:
  bongopay-cli initiate --idempotency-key KEY --currency CODE [flags]
  bongopay-cli get [--server URL] ID

Flags must come before the ID on "get" — Go's flag package stops parsing at the first
non-flag argument, so "get ID --server URL" silently ignores --server.

Run "bongopay-cli initiate -h" or "bongopay-cli get -h" for flags.`)
}

func runInitiate(args []string) {
	fs := flag.NewFlagSet("initiate", flag.ExitOnError)
	server := fs.String("server", "http://localhost:8080", "reference server base URL")
	providerID := fs.String("provider", "SIMULATOR", "Provider.ID")
	amount := fs.Int64("amount", 0, "amount in the currency's minor unit")
	currency := fs.String("currency", "", "ISO 4217 currency code (required)")
	idempotencyKey := fs.String("idempotency-key", "", "idempotency key (required)")
	scenario := fs.String("scenario", "", "simulator scenario name (success, failure, timeout); default: success")
	fs.Parse(args)

	if *currency == "" || *idempotencyKey == "" {
		fmt.Fprintln(os.Stderr, "--currency and --idempotency-key are required")
		fs.Usage()
		os.Exit(2)
	}

	req := payment.PaymentRequest{
		Provider:       payment.Provider{ID: *providerID},
		Amount:         payment.Money{Value: *amount, Currency: payment.Currency{Code: *currency}},
		IdempotencyKey: *idempotencyKey,
	}
	if *scenario != "" {
		raw, err := json.Marshal(map[string]string{"scenario": *scenario})
		if err != nil {
			fatalf("encoding scenario: %v", err)
		}
		req.ProviderOptions = payment.ProviderOptions{"simulator": raw}
	}

	body, err := json.Marshal(req)
	if err != nil {
		fatalf("encoding request: %v", err)
	}

	resp, err := http.Post(*server+"/payments", "application/json", bytes.NewReader(body))
	if err != nil {
		fatalf("request failed: %v", err)
	}
	defer resp.Body.Close()

	printResponse(resp)
}

func runGet(args []string) {
	fs := flag.NewFlagSet("get", flag.ExitOnError)
	server := fs.String("server", "http://localhost:8080", "reference server base URL")
	fs.Parse(args)

	if fs.NArg() < 1 {
		fmt.Fprintln(os.Stderr, "usage: bongopay-cli get [--server URL] ID")
		os.Exit(2)
	}
	id := fs.Arg(0)

	resp, err := http.Get(*server + "/payments/" + id)
	if err != nil {
		fatalf("request failed: %v", err)
	}
	defer resp.Body.Close()

	printResponse(resp)
}

// printResponse pretty-prints resp's JSON body to stdout, and exits 1 if resp was an error
// response — the body (a contracts/openapi/bongopay.yaml Error, on failure) is still printed
// either way, since it's the useful diagnostic.
func printResponse(resp *http.Response) {
	data, err := io.ReadAll(resp.Body)
	if err != nil {
		fatalf("reading response: %v", err)
	}

	var pretty bytes.Buffer
	if err := json.Indent(&pretty, data, "", "  "); err != nil {
		pretty.Write(data) // not JSON — shouldn't happen against this server; print raw
	}
	fmt.Println(pretty.String())

	if resp.StatusCode >= 400 {
		os.Exit(1)
	}
}

func fatalf(format string, args ...any) {
	fmt.Fprintf(os.Stderr, format+"\n", args...)
	os.Exit(1)
}
