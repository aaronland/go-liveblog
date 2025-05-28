package dispatcher

import (
	"context"
	"log/slog"
	"os/exec"
)

type SayDispatcher struct {
	Dispatcher
}

func init() {
	ctx := context.Background()
	err := RegisterDispatcher(ctx, "say", NewSayDispatcher)
	if err != nil {
		panic(err)
	}
}

func NewSayDispatcher(ctx context.Context, uri string) (Dispatcher, error) {
	s := &SayDispatcher{}
	return s, nil
}

func (s *SayDispatcher) Dispatch(ctx context.Context, post string) error {
	slog.Info(post)
	cmd := exec.Command("say", post)
	return cmd.Run()
}
