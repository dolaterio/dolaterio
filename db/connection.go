package db

import (
	"fmt"

	"github.com/dancannon/gorethink"
	"github.com/dolaterio/dolaterio/core"
)

// Connection represents a connection to the database
type Connection struct {
	s            *gorethink.Session
	db           gorethink.Term
	jobsTable    gorethink.Term
	workersTable gorethink.Term
}

// NewConnection initializes and returns a DB connection ready to use
func NewConnection() (*Connection, error) {
	// Open a session to the DB
	s, err := gorethink.Connect(gorethink.ConnectOpts{
		Database: core.Config.RethinkDbDatabase,
		Address: fmt.Sprintf(
			"%v:%v", core.Config.RethinkDbIP, core.Config.RethinkDbPort),
		MaxIdle: 20,
		MaxOpen: 20,
	})

	if err != nil {
		return nil, err
	}
	s.SetMaxOpenConns(5)

	_, err = gorethink.Wait().Run(s)

	if err != nil {
		return nil, err
	}

	db := gorethink.Db(core.Config.RethinkDbDatabase)

	connection := &Connection{
		s:            s,
		db:           db,
		jobsTable:    db.Table("jobs"),
		workersTable: db.Table("workers"),
	}

	return connection, nil
}

func (c *Connection) Close() {
	c.s.Close()
}
