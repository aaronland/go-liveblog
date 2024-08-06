package speaker

import (
	"context"
	"fmt"
	"net/url"
	"sort"
	"strings"

	"github.com/aaronland/go-roster"
)

type Speaker interface {
	ReadPost(context.Context, string) error
}

var speaker_roster roster.Roster

// SpeakerInitializationFunc is a function defined by individual speaker package and used to create
// an instance of that speaker
type SpeakerInitializationFunc func(ctx context.Context, uri string) (Speaker, error)

// RegisterSpeaker registers 'scheme' as a key pointing to 'init_func' in an internal lookup table
// used to create new `Speaker` instances by the `NewSpeaker` method.
func RegisterSpeaker(ctx context.Context, scheme string, init_func SpeakerInitializationFunc) error {

	err := ensureSpeakerRoster()

	if err != nil {
		return err
	}

	return speaker_roster.Register(ctx, scheme, init_func)
}

func ensureSpeakerRoster() error {

	if speaker_roster == nil {

		r, err := roster.NewDefaultRoster()

		if err != nil {
			return err
		}

		speaker_roster = r
	}

	return nil
}

// NewSpeaker returns a new `Speaker` instance configured by 'uri'. The value of 'uri' is parsed
// as a `url.URL` and its scheme is used as the key for a corresponding `SpeakerInitializationFunc`
// function used to instantiate the new `Speaker`. It is assumed that the scheme (and initialization
// function) have been registered by the `RegisterSpeaker` method.
func NewSpeaker(ctx context.Context, uri string) (Speaker, error) {

	u, err := url.Parse(uri)

	if err != nil {
		return nil, err
	}

	scheme := u.Scheme

	i, err := speaker_roster.Driver(ctx, scheme)

	if err != nil {
		return nil, err
	}

	init_func := i.(SpeakerInitializationFunc)
	return init_func(ctx, uri)
}

// SpeakerSchemes returns the list of schemes that have been registered.
func SpeakerSchemes() []string {

	ctx := context.Background()
	schemes := []string{}

	err := ensureSpeakerRoster()

	if err != nil {
		return schemes
	}

	for _, dr := range speaker_roster.Drivers(ctx) {
		scheme := fmt.Sprintf("%s://", strings.ToLower(dr))
		schemes = append(schemes, scheme)
	}

	sort.Strings(schemes)
	return schemes
}
