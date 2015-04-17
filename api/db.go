package api

import (
	"errors"
	"os"

	"github.com/dancannon/gorethink"
)

var (
	// S Represents the database session
	S *gorethink.Session

	// Db Represents the db
	Db gorethink.Term

	// JobTable represents the jobs table
	JobTable gorethink.Term

	errInvalidDbAddress = errors.New(
		"Invalid RethinkDB address. Please make sure $RETHINKDB_ADDRESS is properly set.")
)

// ConnectDb initializes the DB connection
func ConnectDb() error {
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

	S = session

	// Get the db (and create if missing)
	q, err := gorethink.DbList().Run(S)
	if err != nil {
		return err
	}

	var databases []string
	q.All(&databases)

	if !arrContainsString(databases, "dolaterio") {
		_, err = gorethink.DbCreate("dolaterio").RunWrite(S)
		if err != nil {
			return err
		}

	}

	Db = gorethink.Db("dolaterio")

	// Get tables (and create if missing)
	q, err = Db.TableList().Run(S)
	if err != nil {
		return err
	}
	var tables []string
	q.All(&tables)

	if !arrContainsString(tables, "jobs") {
		_, err = Db.TableCreate("jobs").RunWrite(S)
		if err != nil {
			return err
		}
	}
	JobTable = Db.Table("jobs")

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
