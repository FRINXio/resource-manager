package resolver

import (
	"context"

	"github.com/facebookincubator/symphony/pkg/log"
	"github.com/net-auto/resourceManager/ent"
)

type (
	// Config configures resolver.
	Config struct {
		Client *ent.Client
		Logger log.Logger
	}

	// Option allows for managing resolver configuration using functional options.
	Option func(*Resolver)

	Resolver struct {
		logger   log.Logger
		mutation struct{ transactional bool }
	}
)

// New creates a graphql resolver.
func New(cfg Config, opts ...Option) *Resolver {
	r := &Resolver{logger: cfg.Logger}
	r.mutation.transactional = true
	for _, opt := range opts {
		opt(r)
	}
	return r
}

// WithTransaction if set to true, will wraps the mutation with transaction.
func WithTransaction(b bool) Option {
	return func(r *Resolver) {
		r.mutation.transactional = b
	}
}

func (Resolver) ClientFrom(ctx context.Context) *ent.Client {
	client := ent.FromContext(ctx)
	if client == nil {
		panic("no ClientFrom attached to context")
	}
	return client
}
