package lapresse

import (
	"context"
	"log/slog"
	"os/exec"

	"github.com/aaronland/go-liveblog/speaker"
)

type LaPresseSpeaker struct {
	speaker.Speaker
}

func init() {
	ctx := context.Background()
	for _, s := range LaPresseSchemes {
		err := speaker.RegisterSpeaker(ctx, s, NewLaPresseSpeaker)
		if err != nil {
			panic(err)
		}
	}
}

func NewLaPresseSpeaker(ctx context.Context, uri string) (speaker.Speaker, error) {
	s := &LaPresseSpeaker{}
	return s, nil
}

func (s *LaPresseSpeaker) ReadPost(ctx context.Context, post string) error {
	slog.Info(post)
	cmd := exec.Command("say", "-v", "Thomas", post)
	return cmd.Run()
}
