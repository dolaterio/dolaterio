package main

import (
	"fmt"

	"github.com/dancannon/gorethink"
	"github.com/dolaterio/dolaterio/core"
)

func main() {
	s, err := gorethink.Connect(gorethink.ConnectOpts{
		Database: core.Config.RethinkDbDatabase,
		Address: fmt.Sprintf(
			"%v:%v", core.Config.RethinkDbIP, core.Config.RethinkDbPort),
	})
	if err != nil {
		panic(err)
	}
	defer s.Close()

	err = createDb(s)
	if err != nil {
		panic(err)
	}

	err = createTable(s, "jobs")
	if err != nil {
		panic(err)
	}

	err = createTable(s, "workers")
	if err != nil {
		panic(err)
	}
}

func createDb(s *gorethink.Session) error {
	cur, err := gorethink.DbList().Run(s)
	if err != nil {
		return err
	}
	defer cur.Close()

	var databases []string
	cur.All(&databases)

	if !arrContainsString(databases, core.Config.RethinkDbDatabase) {
		_, err = gorethink.DbCreate(core.Config.RethinkDbDatabase).RunWrite(s)
		if err != nil {
			return err
		}
	}
	return nil
}

func createTable(s *gorethink.Session, tableName string) error {
	cur, err := gorethink.Db(core.Config.RethinkDbDatabase).
		TableList().Run(s)
	if err != nil {
		return err
	}
	defer cur.Close()

	var tables []string
	cur.All(&tables)

	if !arrContainsString(tables, tableName) {
		_, err = gorethink.Db(core.Config.RethinkDbDatabase).
			TableCreate(tableName).RunWrite(s)
		if err != nil {
			return err
		}
	}
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
