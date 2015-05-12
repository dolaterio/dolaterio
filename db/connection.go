package db

import "github.com/dancannon/gorethink"

// Connection represents a connection to the database
type Connection struct {
	s         *gorethink.Session
	db        gorethink.Term
	jobsTable gorethink.Term
}

// NewConnection initializes and returns a DB connection ready to use
func NewConnection(conf *ConnectionConfig) (*Connection, error) {
	conf.defaults()
	if err := conf.errors(); err != nil {
		return nil, err
	}

	// Open a session to the DB
	s, err := gorethink.Connect(gorethink.ConnectOpts{
		Address:  conf.RethinkDbAddress,
		Database: "dolaterio",
		MaxIdle:  10,
		MaxOpen:  10,
	})

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

	return connection, nil
}

func arrContainsString(arr []string, val string) bool {
	for _, it := range arr {
		if it == val {
			return true
		}
	}
	return false
}
