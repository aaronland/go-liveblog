package guardian

import (
	"context"
	_ "log/slog"

	"github.com/aaronland/go-liveblog/dispatcher"
)

func init() {
	ctx := context.Background()
	for _, s := range GuardianSchemes {
		err := dispatcher.RegisterDispatcher(ctx, s, dispatcher.NewSayDispatcher)
		if err != nil {
			panic(err)
		}
	}
}
