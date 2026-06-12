//go:build darwin

package dispatcher

import (
	"context"
	"os/exec"
)

// SayDispatcher implements the Dispatcher interface by utilizing the
// MacOS "say" command to output text as speech.
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

// NewSayDispatcher initializes a new SayDispatcher instance.
func NewSayDispatcher(ctx context.Context, uri string) (Dispatcher, error) {
	s := &SayDispatcher{}
	return s, nil
}

// Dispatch executes the system "say" command with the provided string as the argument.
func (s *SayDispatcher) Dispatch(ctx context.Context, post string) error {
	cmd := exec.Command("say", post)
	return cmd.Run()
}
