package api

import "github.com/dancannon/gorethink"

// Session Represents the database session
var Session *gorethink.Session

// Db Represents the db
var Db gorethink.Term

// JobTable represents the jobs table
var JobTable gorethink.Term

func connectDb() *gorethink.Session {
	session, err := gorethink.Connect(gorethink.ConnectOpts{
		Address:  "d.lo:28015",
		Database: "test",
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
