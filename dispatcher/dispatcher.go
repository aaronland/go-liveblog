package dispatcher

import (
	"context"
	"fmt"
	"net/url"
	"sort"
	"strings"

	"github.com/aaronland/go-roster"
)

type Dispatcher interface {
	Dispatch(context.Context, string) error
}

var dispatcher_roster roster.Roster

// DispatcherInitializationFunc is a function defined by individual dispatcher package and used to create
// an instance of that dispatcher
type DispatcherInitializationFunc func(ctx context.Context, uri string) (Dispatcher, error)

// RegisterDispatcher registers 'scheme' as a key pointing to 'init_func' in an internal lookup table
// used to create new `Dispatcher` instances by the `NewDispatcher` method.
func RegisterDispatcher(ctx context.Context, scheme string, init_func DispatcherInitializationFunc) error {

	err := ensureDispatcherRoster()

	if err != nil {
		return err
	}

	return dispatcher_roster.Register(ctx, scheme, init_func)
}

func ensureDispatcherRoster() error {

	if dispatcher_roster == nil {

		r, err := roster.NewDefaultRoster()

		if err != nil {
			return err
		}

		dispatcher_roster = r
	}

	return nil
}

// NewDispatcher returns a new `Dispatcher` instance configured by 'uri'. The value of 'uri' is parsed
// as a `url.URL` and its scheme is used as the key for a corresponding `DispatcherInitializationFunc`
// function used to instantiate the new `Dispatcher`. It is assumed that the scheme (and initialization
// function) have been registered by the `RegisterDispatcher` method.
func NewDispatcher(ctx context.Context, uri string) (Dispatcher, error) {

	u, err := url.Parse(uri)

	if err != nil {
		return nil, err
	}

	scheme := u.Scheme

	i, err := dispatcher_roster.Driver(ctx, scheme)

	if err != nil {
		return nil, err
	}

	init_func := i.(DispatcherInitializationFunc)
	return init_func(ctx, uri)
}

// DispatcherSchemes returns the list of schemes that have been registered.
func DispatcherSchemes() []string {

	ctx := context.Background()
	schemes := []string{}

	err := ensureDispatcherRoster()

	if err != nil {
		return schemes
	}

	for _, dr := range dispatcher_roster.Drivers(ctx) {
		scheme := fmt.Sprintf("%s://", strings.ToLower(dr))
		schemes = append(schemes, scheme)
	}

	sort.Strings(schemes)
	return schemes
}
