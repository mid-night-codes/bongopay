module github.com/mid-night-codes/bongopay/examples/go-quickstart

go 1.27.1

require github.com/mid-night-codes/bongopay/sdks/go v0.0.0

// sdks/go isn't published anywhere — this always points at the copy in this same repo.
replace github.com/mid-night-codes/bongopay/sdks/go => ../../sdks/go
