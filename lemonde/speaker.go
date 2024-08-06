package lemonde

import (
	"context"
	"log/slog"
	"os/exec"

	"github.com/aaronland/go-liveblog/speaker"
)

type LeMondeSpeaker struct {
	speaker.Speaker
}

func init() {
	ctx := context.Background()
	for _, s := range LeMondeSchemes {
		err := speaker.RegisterSpeaker(ctx, s, NewLeMondeSpeaker)
		if err != nil {
			panic(err)
		}
	}
}

func NewLeMondeSpeaker(ctx context.Context, uri string) (speaker.Speaker, error) {
	s := &LeMondeSpeaker{}
	return s, nil
}

func (s *LeMondeSpeaker) ReadPost(ctx context.Context, post string) error {
	slog.Info(post)
	cmd := exec.Command("say", "-v", "Thomas", post)
	return cmd.Run()
}
