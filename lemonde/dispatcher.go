package lemonde

import (
	"context"
	"log/slog"
	"os/exec"

	"github.com/aaronland/go-liveblog/dispatcher"
)

type LeMondeDispatcher struct {
	dispatcher.Dispatcher
}

func init() {
	ctx := context.Background()
	for _, s := range LeMondeSchemes {
		err := dispatcher.RegisterDispatcher(ctx, s, NewLeMondeDispatcher)
		if err != nil {
			panic(err)
		}
	}
}

func NewLeMondeDispatcher(ctx context.Context, uri string) (dispatcher.Dispatcher, error) {
	s := &LeMondeDispatcher{}
	return s, nil
}

func (s *LeMondeDispatcher) Dispatch(ctx context.Context, post string) error {
	slog.Info(post)
	cmd := exec.Command("say", "-v", "Thomas", post)
	return cmd.Run()
}
