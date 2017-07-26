package selector

import (
	"math/rand"
	"time"

	"registry"
)

type randomSelector struct {
	so Options
}

func init() {
	rand.Seed(time.Now().Unix())
}

func (r *randomSelector) Init(opts ...Option) error {
	for _, o := range opts {
		o(&r.so)
	}
	return nil
}

func (r *randomSelector) Options() Options {
	return r.so
}

func (r *randomSelector) Select(service string, opts ...SelectOption) (Next, error) {
	sopts := SelectOptions{
		Strategy: r.so.Strategy,
	}

	for _, opt := range opts {
		opt(&sopts)
	}

	// get the service
	services, err := r.so.Registry.GetService(service)
	if err != nil {
		return nil, err
	}

	// apply the filters
	for _, filter := range sopts.Filters {
		services = filter(services)
	}

	// if there's nothing left, return
	if len(services) == 0 {
		return nil, ErrNoneAvailable
	}

	return sopts.Strategy(services), nil
}

func (r *randomSelector) Mark(service string, node *registry.Node, err error) {
	return
}

func (r *randomSelector) Reset(service string) {
	return
}

func (r *randomSelector) Close() error {
	return nil
}

func (r *randomSelector) String() string {
	return "random"
}

func newRandomSelector(opts ...Option) Selector {
	sopts := Options{
		Strategy: Random,
	}

	for _, opt := range opts {
		opt(&sopts)
	}

	if sopts.Registry == nil {
		sopts.Registry = registry.DefaultRegistry
	}

	return &randomSelector{
		so: sopts,
	}
}
