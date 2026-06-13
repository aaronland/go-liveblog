package follow

import (
	"context"
	"flag"
	"fmt"
	"log/slog"
	"net/http"
	net_url "net/url"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/aaronland/go-liveblog/app/follow/www"
	"github.com/aaronland/go-liveblog/dispatcher"
	"github.com/aaronland/go-liveblog/parser"
	"github.com/sfomuseum/go-flags/flagset"
	"github.com/sfomuseum/go-pubsub/publisher"
	"github.com/sfomuseum/go-pubsub/subscriber"
	"github.com/sfomuseum/go-www-show/v2"
	"github.com/whosonfirst/go-pubssed/broker"
)

const WWW_DISPATCHER string = "web://"

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

	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	sig_ch := make(chan os.Signal, 1)
	signal.Notify(sig_ch, os.Interrupt, syscall.SIGTERM, syscall.SIGHUP)

	go func() {
		for s := range sig_ch {
			switch s {
			case os.Interrupt, syscall.SIGTERM:
				slog.Info("Shutdown signal received.")
				cancel()
				os.Exit(0)
			}
		}
	}()

	urls := fs.Args()

	var dp dispatcher.Dispatcher

	dp_u, err := net_url.Parse(dispatcher_uri)

	if err != nil {
		return fmt.Errorf("Failed to parse dispatcher URI, %w", err)
	}

	switch dp_u.Scheme {
	case "webview": // read from constants in go-www-show/v2...
		return fmt.Errorf("webview:// not supported at this time")
	case "web": // read from constants in go-www-show/v2...

		dp_ch := make(chan string)

		pub, err := publisher.NewChannelPublisherWithChannel(ctx, dp_ch)

		if err != nil {
			return fmt.Errorf("Failed to create channel publisher, %w", err)
		}

		sub, err := subscriber.NewChannelSubscriberWithChannel(ctx, dp_ch)

		if err != nil {
			return fmt.Errorf("Failed to create channel subscribed, %w", err)
		}

		defer sub.Close()

		d, err := dispatcher.NewPubSubDispatcherWithPublisher(ctx, pub)

		if err != nil {
			return fmt.Errorf("Failed to create pubsub dispatcher, %w", err)
		}

		dp = d

		mux := http.NewServeMux()

		brkr, err := broker.NewBroker()

		if err != nil {
			return fmt.Errorf("Failed to create pubsub-sse broker, %w", err)
		}

		sse_handler, err := brkr.HandlerFunc()

		if err != nil {
			return fmt.Errorf("Failed to create pubsub-sse handler, %w", err)
		}

		brkr.Start(ctx, sub)

		mux.HandleFunc("/sse", sse_handler)

		http_fs := http.FS(www.FS)
		index_handler := http.FileServer(http_fs)

		mux.Handle("/", index_handler)

		browser, err := show.NewBrowser(ctx, dispatcher_uri)

		if err != nil {
			return fmt.Errorf("Failed to create new browser, %w", err)
		}

		show_opts := &show.RunOptions{
			Browser:            browser,
			Mux:                mux,
			EnsureServiceDelay: 200 * time.Millisecond,
		}

		go func() {

			err := show.RunWithOptions(ctx, show_opts)

			if err != nil {
				panic(err)
			}
		}()

	default:

		d, err := dispatcher.NewDispatcher(ctx, dispatcher_uri)

		if err != nil {
			return fmt.Errorf("Failed to create new dispatcher, %w", err)
		}

		dp = d
	}

	cache := new(sync.Map)
	mu := new(sync.RWMutex)

	process(ctx, dp, cache, mu, read_all, urls...)

	ticker := time.NewTicker(time.Duration(delay) * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			break
		case <-ticker.C:
			slog.Debug("Process URIs")
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

	logger := slog.Default()
	logger = logger.With("url", url)

	logger.Debug("Handle posts", "read", read)

	u, err := net_url.Parse(url)

	if err != nil {
		logger.Error("Failed to parse URL", "error", err)
		return
	}

	p_uri := fmt.Sprintf("%s://", u.Host)
	p, err := parser.NewParser(ctx, p_uri)

	if err != nil {
		logger.Error("Failed to derive new parser", "uri", p_uri, "error", err)
	}

	title, posts, err := p.GetPosts(ctx, url)

	if err != nil {
		logger.Error("Failed to retrieve posts", "url", url, "error", err)
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

			err := dp.Dispatch(ctx, p)

			if err != nil {
				logger.Error("Failed to dispatch post", "error", err)
			}
		}
	}
}
