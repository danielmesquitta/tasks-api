package mysqlrepo

import (
	"database/sql"

	_ "github.com/go-sql-driver/mysql" // mysql driver

	"github.com/danielmesquitta/tasks-api/internal/config"
	"github.com/danielmesquitta/tasks-api/internal/provider/db/mysqldb"
)

func NewMySQLDBConn(
	env *config.Env,
) *mysqldb.Queries {
	dbConn, err := sql.Open(
		"mysql",
		env.DBConnection,
	)
	if err != nil {
		panic(err)
	}

	err = dbConn.Ping()
	if err != nil {
		panic(err)
	}

	return mysqldb.New(dbConn)
}
