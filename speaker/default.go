package speaker

import (
	"context"
	"log/slog"
	"os/exec"
)

type DefaultSpeaker struct {
	Speaker
}

func init() {
	ctx := context.Background()
	err := RegisterSpeaker(ctx, "default", NewDefaultSpeaker)
	if err != nil {
		panic(err)
	}
}

func NewDefaultSpeaker(ctx context.Context, uri string) (Speaker, error) {
	s := &DefaultSpeaker{}
	return s, nil
}

func (s *DefaultSpeaker) ReadPost(ctx context.Context, post string) error {
	slog.Info(post)
	cmd := exec.Command("say", post)
	return cmd.Run()
}
