package dispatcher

import (
	"context"
	"fmt"
	"strings"

	"github.com/sfomuseum/go-pubsub/publisher"
)

// PubSubDispatcher implements the Dispatcher interface for a PubSub-based
// messaging system using the `sfomuseum/go-pubsub` package.
type PubSubDispatcher struct {
	Dispatcher
	publisher publisher.Publisher
}

func init() {

	ctx := context.Background()

	for _, scheme := range publisher.PublisherSchemes() {

		scheme = strings.Replace(scheme, "://", "", 1)

		err := RegisterDispatcher(ctx, scheme, NewPubSubDispatcher)

		if err != nil {
			panic(err)
		}
	}
}

// NewPubSubDispatcher creates a new PubSubDispatcher by initializing the
// underlying publisher using the provided URI which is expected to be a
// registered `sfomuseum/go-pubsub/publisher` URI.
func NewPubSubDispatcher(ctx context.Context, uri string) (Dispatcher, error) {

	p, err := publisher.NewPublisher(ctx, uri)

	if err != nil {
		return nil, fmt.Errorf("Failed to create new publisher, %w", err)
	}

	return NewPubSubDispatcherWithPublisher(ctx, p)
}

// NewPubSubDispatcherWithPublisher is a helper constructor that takes an
// already initialized publisher.
func NewPubSubDispatcherWithPublisher(ctx context.Context, p publisher.Publisher) (Dispatcher, error) {

	s := &PubSubDispatcher{
		publisher: p,
	}

	return s, nil
}

// Dispatch sends the payload to the PubSub publisher.
func (s *PubSubDispatcher) Dispatch(ctx context.Context, post string) error {
	return s.publisher.Publish(ctx, post)
}
