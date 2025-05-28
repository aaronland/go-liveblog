package dispatcher

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/sfomuseum/go-pubsub/publisher"
)

type PubSubDispatcher struct {
	Dispatcher
	publisher publisher.Publisher
}

func init() {
	ctx := context.Background()

	for _, scheme := range publisher.PublisherSchemes() {

		err := RegisterDispatcher(ctx, scheme, NewPubSubDispatcher)
		if err != nil {
			panic(err)
		}
	}
}

func NewPubSubDispatcher(ctx context.Context, uri string) (Dispatcher, error) {

	p, err := publisher.NewPublisher(ctx, uri)

	if err != nil {
		return nil, fmt.Errorf("Failed to create new publisher, %w", err)
	}

	return NewPubSubDispatcherWithPublisher(ctx, p)
}

func NewPubSubDispatcherWithPublisher(ctx context.Context, p publisher.Publisher) (Dispatcher, error) {

	s := &PubSubDispatcher{
		publisher: p,
	}

	return s, nil
}

func (s *PubSubDispatcher) Dispatch(ctx context.Context, post string) error {
	slog.Info(post)
	return s.publisher.Publish(ctx, post)
}
