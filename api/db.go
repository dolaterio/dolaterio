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

	INVALID_DB_ADDRESS = errors.New(
		"Invalid RethinkDB address. Please make sure $RETHINKDB_ADDRESS is properly set.")
)

// ConnectDb initializes the DB connection
func ConnectDb() error {
	address := os.Getenv("RETHINKDB_ADDRESS")
	if address == "" {
		return INVALID_DB_ADDRESS
	}
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

	gorethink.DbCreate("dolaterio").RunWrite(S)
	Db = gorethink.Db("dolaterio")

	Db.TableCreate("jobs").RunWrite(S)
	JobTable = Db.Table("jobs")

	return nil
}
