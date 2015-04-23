package db

import (
	"errors"
	"os"

	"github.com/dancannon/gorethink"
)

var (
	s         *gorethink.Session
	db        gorethink.Term
	jobsTable gorethink.Term

	errInvalidDbAddress = errors.New(
		"Invalid RethinkDB address. Please make sure $RETHINKDB_ADDRESS is properly set.")
)

// Connect initializes the DB connection
func Connect() error {
	if s != nil {
		// It's already connected
		return nil
	}

	address := os.Getenv("RETHINKDB_ADDRESS")
	if address == "" {
		return errInvalidDbAddress
	}

	// Open a session to the DB
	session, err := gorethink.Connect(gorethink.ConnectOpts{
		Address:  os.Getenv("RETHINKDB_ADDRESS"),
		Database: "dolaterio",
		MaxIdle:  10,
		MaxOpen:  10,
	})

	if err != nil {
		return err
	}

	session.SetMaxOpenConns(5)

	s = session

	// Get the db (and create if missing)
	q, err := gorethink.DbList().Run(s)
	if err != nil {
		return err
	}

	var databases []string
	q.All(&databases)

	if !arrContainsString(databases, "dolaterio") {
		_, err = gorethink.DbCreate("dolaterio").RunWrite(s)
		if err != nil {
			return err
		}

	}

	db = gorethink.Db("dolaterio")

	// Get tables (and create if missing)
	q, err = db.TableList().Run(s)
	if err != nil {
		return err
	}
	var tables []string
	q.All(&tables)

	if !arrContainsString(tables, "jobs") {
		_, err = db.TableCreate("jobs").RunWrite(s)
		if err != nil {
			return err
		}
	}
	jobsTable = db.Table("jobs")

	return nil
}

func arrContainsString(arr []string, val string) bool {
	for _, it := range arr {
		if it == val {
			return true
		}
	}
	return false
}
