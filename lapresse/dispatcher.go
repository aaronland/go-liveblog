package lapresse

import (
	"context"
	"log/slog"
	"os/exec"

	"github.com/aaronland/go-liveblog/dispatcher"
)

type LaPresseDispatcher struct {
	dispatcher.Dispatcher
}

func init() {
	ctx := context.Background()
	for _, s := range LaPresseSchemes {
		err := dispatcher.RegisterDispatcher(ctx, s, NewLaPresseDispatcher)
		if err != nil {
			panic(err)
		}
	}
}

func NewLaPresseDispatcher(ctx context.Context, uri string) (dispatcher.Dispatcher, error) {
	s := &LaPresseDispatcher{}
	return s, nil
}

func (s *LaPresseDispatcher) Dispatch(ctx context.Context, post string) error {
	slog.Info(post)
	cmd := exec.Command("say", "-v", "Thomas", post)
	return cmd.Run()
}
