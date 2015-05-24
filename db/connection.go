package db

import (
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
		Address: core.Config.RethinkDbAddress,
		MaxIdle: 20,
		MaxOpen: 20,
	})

	if err != nil {
		return nil, err
	}

	_, err = gorethink.Wait().Run(s)

	if err != nil {
		return nil, err
	}

	connection := &Connection{
		s: s,
	}

	s.SetMaxOpenConns(5)

	// Get the db (and create if missing)
	q, err := gorethink.DbList().Run(s)
	if err != nil {
		return nil, err
	}

	var databases []string
	q.All(&databases)

	if !arrContainsString(databases, "dolaterio") {
		_, err = gorethink.DbCreate("dolaterio").RunWrite(s)
		if err != nil {
			return nil, err
		}

	}

	connection.db = gorethink.Db("dolaterio")

	_, err = connection.db.Wait().Run(s)
	if err != nil {
		return nil, err
	}

	// Get tables (and create if missing)
	q, err = connection.db.TableList().Run(s)
	if err != nil {
		return nil, err
	}
	var tables []string
	q.All(&tables)

	if !arrContainsString(tables, "jobs") {
		_, err = connection.db.TableCreate("jobs").RunWrite(s)
		if err != nil {
			return nil, err
		}
	}
	connection.jobsTable = connection.db.Table("jobs")
	_, err = connection.jobsTable.Wait().Run(s)
	if err != nil {
		return nil, err
	}

	if !arrContainsString(tables, "workers") {
		_, err = connection.db.TableCreate("workers").RunWrite(s)
		if err != nil {
			return nil, err
		}
	}
	connection.workersTable = connection.db.Table("workers")
	_, err = connection.workersTable.Wait().Run(s)
	if err != nil {
		return nil, err
	}

	return connection, nil
}

func (c *Connection) Close() {
	c.s.Close()
}

func arrContainsString(arr []string, val string) bool {
	for _, it := range arr {
		if it == val {
			return true
		}
	}
	return false
}
