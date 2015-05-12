package db

import (
	"errors"
	"os"
)

var (
	errInvalidDbAddress = errors.New(
		"Invalid RethinkDB address. Please make sure $RETHINKDB_ADDRESS is properly set.")
)

// ConnectionConfig reprensents the config values to stablish a db connection
type ConnectionConfig struct {
	RethinkDbAddress string
}

func (c *ConnectionConfig) defaults() {
	if c.RethinkDbAddress == "" {
		c.RethinkDbAddress = os.Getenv("RETHINKDB_ADDRESS")
	}
}

func (c *ConnectionConfig) errors() error {
	if c.RethinkDbAddress == "" {
		return errInvalidDbAddress
	}
	return nil
}
