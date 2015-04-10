package api

import (
	"fmt"

	"github.com/dancannon/gorethink"
)

var (
	// S Represents the database session
	S *gorethink.Session

	// Db Represents the db
	Db gorethink.Term

	// JobTable represents the jobs table
	JobTable gorethink.Term
)

// ConnectDb initializes the DB connection
func ConnectDb(rdbHost, rdbPort string) error {
	rdbAddress := fmt.Sprintf("%v:%v", rdbHost, rdbPort)
	session, err := gorethink.Connect(gorethink.ConnectOpts{
		Address:  rdbAddress,
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
