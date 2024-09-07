package mysqlrepo

import (
	"context"
	"database/sql"

	_ "github.com/go-sql-driver/mysql" // mysql driver

	"github.com/danielmesquitta/tasks-api/internal/config"
	"github.com/danielmesquitta/tasks-api/internal/pkg/transactioner"
	"github.com/danielmesquitta/tasks-api/internal/provider/db/mysqldb"
)

func NewMySQLDBConn(
	env *config.Env,
) *sql.DB {
	dbConn, err := sql.Open(
		"mysql",
		env.DBConnection,
	)
	if err != nil {
		panic(err)
	}

	if err := dbConn.Ping(); err != nil {
		panic(err)
	}

	return dbConn
}

type Queries struct {
	*mysqldb.Queries
}

func NewMySQLQueries(
	dbConn *sql.DB,
) *Queries {
	return &Queries{mysqldb.New(dbConn)}
}

func (q *Queries) getDBorTX(
	ctx context.Context,
) *Queries {
	tx, ok := ctx.Value(transactioner.CtxTxKey).(*sql.Tx)
	if ok {
		return &Queries{q.WithTx(tx)}
	}
	return q
}
