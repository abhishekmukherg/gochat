package gchatdb

import (
	"database/sql"
	"log"
)

const latestVersion = 1

type DbConnection struct {
	*sql.DB
}

func Open(driverName, dataSourceName string) (*DbConnection, error) {
	db, err := sql.Open(driverName, dataSourceName)
	return &DbConnection{db}, err
}

func (db *DbConnection) getVersion() (int, error) {
	var cur_version int
	err := db.QueryRow("SELECT cur_version FROM versions").Scan(&cur_version)
	switch {
	case err == sql.ErrNoRows:
		log.Fatalf("Got no rows back, that shouldn't happen: %q\n", err)
		return 0, err
	case err != nil:
		log.Printf("Got an error, assuming we just didn't have table: %q\n", err)
		return 0, nil
	default:
		return cur_version, nil
	}
}

func (db *DbConnection) Upgrade() error {
	version, err := db.getVersion()
	if err != nil {
		return err
	}

	if version == latestVersion {
		return nil
	}

	err = db.doUpgrade(version)
	if err != nil {
		return err
	}

	err = db.updateVersion(latestVersion)
	if err != nil {
		return err
	}

	return nil
}

func (db *DbConnection) doUpgrade(version int) error {
	sqlStmts := []string{
		"create table versions (cur_version integer)",
		"create table users (id integer not null primary key, name text unique not null)",
	}

	for _, sqlStmt := range sqlStmts {
		_, err := db.Exec(sqlStmt)
		if err != nil {
			return err
		}
	}

	return nil
}

func (db *DbConnection) updateVersion(version int) error {
	sqlStmt := `
	DELETE FROM versions;
	INSERT INTO versions VALUES (?);
	`
	_, err := db.Exec(sqlStmt, version)
	if err != nil {
		return err
	}
	return nil
}

