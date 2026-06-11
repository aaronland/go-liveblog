package follow

import (
	"context"
	"flag"
	"fmt"
	"log/slog"
	"net/http"
	net_url "net/url"
	"sync"
	"time"

	"github.com/aaronland/go-liveblog/dispatcher"
	"github.com/aaronland/go-liveblog/parser"
	"github.com/sfomuseum/go-flags/flagset"
	"github.com/sfomuseum/go-pubsub/subscriber"
	"github.com/whosonfirst/go-pubssed/broker"
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

	urls := fs.Args()

	dp, err := dispatcher.NewDispatcher(ctx, dispatcher_uri)

	if err != nil {
		return err
	}

	if www {

		sub, err := subscriber.NewSubscriber(ctx, "redis://?host=localhost&port=6379&channel=pubssed")

		if err != nil {
			return err
		}

		brkr, err := broker.NewBroker()

		if err != nil {
			return err
		}

		http_handler, err := brkr.HandlerFunc()

		if err != nil {
			return err
		}

		brkr.Start(ctx, sub)

		mux := http.NewServeMux()
		mux.HandleFunc("/", http_handler)

		go http.ListenAndServe("localhost:8080", mux)
	}

	cache := new(sync.Map)
	mu := new(sync.RWMutex)

	process(ctx, dp, cache, mu, read_all, urls...)

	ticker := time.NewTicker(time.Duration(delay) * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			process(ctx, dp, cache, mu, true, urls...)
		}
	}

	return nil
}

func process(ctx context.Context, dp dispatcher.Dispatcher, cache *sync.Map, mu *sync.RWMutex, read bool, urls ...string) {

	read_title := false

	if len(urls) > 1 {
		read_title = true
	}

	for _, url := range urls {
		go handle_posts(ctx, dp, cache, mu, read, read_title, url)
	}

}

func handle_posts(ctx context.Context, dp dispatcher.Dispatcher, cache *sync.Map, mu *sync.RWMutex, read bool, read_title bool, url string) {

	slog.Debug("Handle posts", "url", url, "read", read)

	u, err := net_url.Parse(url)

	if err != nil {
		slog.Error("Failed to parse URL", "error", err)
		return
	}

	p_uri := fmt.Sprintf("%s://", u.Host)
	p, err := parser.NewParser(ctx, p_uri)

	if err != nil {
		slog.Error("Failed to derive new parser", "uri", p_uri, "error", err)
	}

	title, posts, err := p.GetPosts(ctx, url)

	if err != nil {
		slog.Error("Failed to retrieve posts", "url", url, "error", err)
		return
	}

	title_read := true

	mu.Lock()
	defer mu.Unlock()

	for _, p := range posts {

		_, exists := cache.LoadOrStore(p, true)

		if exists {
			continue
		}

		if read {

			if read_title && !title_read {
				dp.Dispatch(ctx, title)
				title_read = false
			}

			dp.Dispatch(ctx, p)
		}
	}
}
