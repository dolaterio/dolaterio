package runner

import "github.com/dolaterio/dolaterio/db"

var (
	dbConnection *db.Connection
)

func init() {
	dbConnection, _ = db.NewConnection(&db.ConnectionConfig{})
}
