package gchatdb

import (
	"database/sql"
)

type DbConnection *sql.DB

func Open(driverName, dataSourceName string) (*DbConnection, error) {
	return sql.Open(driverName, dataSourceName)
}
