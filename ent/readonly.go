// Code generated by entc, DO NOT EDIT.

package ent

import (
	"context"

	"entgo.io/ent/dialect"
	"github.com/net-auto/resourceManager/ent/migrate"
)

// ReadOnly returns a new readonly-client.
//
//	client := client.ReadOnly()
func (c *Client) ReadOnly() *Client {
	cfg := config{driver: &readonly{Driver: c.driver}, log: c.log}
	return &Client{
		config:             cfg,
		Schema:             migrate.NewSchema(cfg.driver),
		AllocationStrategy: NewAllocationStrategyClient(cfg),
		PoolProperties:     NewPoolPropertiesClient(cfg),
		Property:           NewPropertyClient(cfg),
		PropertyType:       NewPropertyTypeClient(cfg),
		Resource:           NewResourceClient(cfg),
		ResourcePool:       NewResourcePoolClient(cfg),
		ResourceType:       NewResourceTypeClient(cfg),
		Tag:                NewTagClient(cfg),
	}
}

// ErrReadOnly returns when a readonly user tries to execute a write operation.
var ErrReadOnly = &PermissionError{cause: "permission denied: read-only user"}

// PermissionError represents a permission denied error.
type PermissionError struct {
	cause string
}

func (e PermissionError) Error() string { return e.cause }

type readonly struct {
	dialect.Driver
}

func (r *readonly) Exec(context.Context, string, interface{}, interface{}) error {
	return ErrReadOnly
}

func (r *readonly) Tx(context.Context) (dialect.Tx, error) {
	return nil, ErrReadOnly
}
