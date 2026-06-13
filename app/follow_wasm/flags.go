package follow_wasm

import (
	"flag"
	"fmt"
	"os"

	"github.com/sfomuseum/go-flags/flagset"
)

var verbose bool

var server_uri string

func DefaultFlagSet() *flag.FlagSet {

	fs := flagset.NewFlagSet("follow")

	fs.BoolVar(&verbose, "verbose", false, "Enable verbose (debug) logging.")
	fs.StringVar(&server_uri, "server-uri", "http://localhost:8080", "...")

	fs.Usage = func() {
		fmt.Fprintf(os.Stderr, "\n")
		fmt.Fprintf(os.Stderr, "Usage:\n\t %s [options]\n", os.Args[0])
		fs.PrintDefaults()
	}

	return fs
}
