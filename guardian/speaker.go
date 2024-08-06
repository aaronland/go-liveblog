package guardian

import (
	"context"
	_ "log/slog"

	"github.com/aaronland/go-liveblog/speaker"
)

func init() {
	ctx := context.Background()
	for _, s := range GuardianSchemes {
		err := speaker.RegisterSpeaker(ctx, s, speaker.NewDefaultSpeaker)
		if err != nil {
			panic(err)
		}
	}
}
