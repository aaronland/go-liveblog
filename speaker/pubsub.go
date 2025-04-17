package speaker

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/sfomuseum/go-pubsub/publisher"
)

type PubSubSpeaker struct {
	Speaker
	publisher publisher.Publisher
}

func init() {
	ctx := context.Background()

	for _, scheme := range publisher.PublisherSchemes() {

		err := RegisterSpeaker(ctx, scheme, NewPubSubSpeaker)
		if err != nil {
			panic(err)
		}
	}
}

func NewPubSubSpeaker(ctx context.Context, uri string) (Speaker, error) {

	p, err := publisher.NewPublisher(ctx, uri)

	if err != nil {
		return nil, fmt.Errorf("Failed to create new publisher, %w", err)
	}

	return NewPubSubSpeakerWithPublisher(ctx, p)
}

func NewPubSubSpeakerWithPublisher(ctx context.Context, p publisher.Publisher) (Speaker, error) {

	s := &PubSubSpeaker{
		publisher: p,
	}

	return s, nil
}

func (s *PubSubSpeaker) ReadPost(ctx context.Context, post string) error {
	slog.Info(post)
	return s.publisher.Publish(ctx, post)
}
