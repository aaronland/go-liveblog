package follow_wasm

import (
	"context"
	"flag"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/aaronland/go-http/v4/server"
	"github.com/aaronland/go-liveblog/app/follow_wasm/www"
	"github.com/sfomuseum/go-flags/flagset"
)

func Run(ctx context.Context) error {
	fs := DefaultFlagSet()
	return RunWithFlagSet(ctx, fs)
}

func RunWithFlagSet(ctx context.Context, fs *flag.FlagSet) error {

	flagset.Parse(fs)

	if verbose {
		slog.SetLogLoggerLevel(slog.LevelDebug)
		slog.Debug("Verbose logging enabled")
	}

	mux := http.NewServeMux()

	http_fs := http.FS(www.FS)
	index_handler := http.FileServer(http_fs)

	mux.Handle("/", index_handler)

	s, err := server.NewServer(ctx, server_uri)

	if err != nil {
		return fmt.Errorf("Failed to create server, %w", err)
	}

	slog.Info("Listening for requests", "address", s.Address())
	return s.ListenAndServe(ctx, mux)
}
