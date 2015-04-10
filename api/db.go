package api

import (
	"os"

	"github.com/dancannon/gorethink"
)

// Session Represents the database session
var Session *gorethink.Session

// Db Represents the db
var Db gorethink.Term

// JobTable represents the jobs table
var JobTable gorethink.Term

func connectDb() *gorethink.Session {
	rdbHost := os.Getenv("RETHINKDB_PORT_28015_TCP_ADDR")
	if rdbHost == "" {
		rdbHost = "d.lo"
	}
	rdbPort := os.Getenv("RETHINKDB_PORT_28015_TCP_PORT")
	if rdbPort == "" {
		rdbPort = "28015"
	}

	session, err := gorethink.Connect(gorethink.ConnectOpts{
		Address:  rdbHost + ":" + rdbPort,
		Database: "dolaterio",
		MaxIdle:  10,
		MaxOpen:  10,
	})

	if err != nil {
		panic(err)
	}

	session.SetMaxOpenConns(5)

	Session = session

	gorethink.DbCreate("dolaterio").RunWrite(Session)
	Db = gorethink.Db("dolaterio")

	Db.TableCreate("jobs").RunWrite(Session)
	JobTable = Db.Table("jobs")

	return session
}
