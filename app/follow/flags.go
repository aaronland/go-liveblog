package follow

import (
	"flag"
	"fmt"
	"os"

	"github.com/sfomuseum/go-flags/flagset"
)

var delay int
var read_all bool
var verbose bool

var www_ui bool
var dispatcher_uri string

func DefaultFlagSet() *flag.FlagSet {

	fs := flagset.NewFlagSet("follow")

	fs.IntVar(&delay, "delay", 30, "The number of seconds to wait before fetching new updates")
	fs.BoolVar(&read_all, "read-all", false, "If true read all previous posts (written before following has begun)")
	fs.BoolVar(&verbose, "verbose", false, "Enable verbose (debug) logging.")
	fs.StringVar(&dispatcher_uri, "dispatcher-uri", "say://", "...")
	fs.BoolVar(&www_ui, "www", false, "...")

	fs.Usage = func() {
		fmt.Fprintf(os.Stderr, "Parse one or more \"live blog\" URLs and read them aloud.\n")
		fmt.Fprintf(os.Stderr, "Usage:\n\t %s [options] url(N) url(N)\n", os.Args[0])
		fs.PrintDefaults()
	}

	return fs
}
