package gchatdb

import (
	"database/sql"
	"log"
)

var (
	dbQueries = []string{
		"CREATE TABLE versions (cur_version INTEGER)",
		`CREATE TABLE users (id INTEGER NOT NULL PRIMARY KEY,
					name TEXT UNIQUE NOT NULL,
					password BLOB,
				        passwordVersion INTEGER)`,
		`CREATE TABLE chats (id INTEGER NOT NULL PRIMARY KEY,
					title TEXT UNIQUE NOT NULL)`,
		`CREATE TABLE chat_users (chat_id INTEGER NOT NULL,
				          user_id INTEGER NOT NULL)`,
		`CREATE UNIQUE INDEX chat_users_idx ON chat_users (chat_id, user_id)`,

		`CREATE TABLE messages (id INTEGER NOT NULL PRIMARY KEY,
					user_id INTEGER NOT NULL,
				        chat_id INTEGER NOT NULL,
				        timestamp INTEGER NOT NULL,
				        text TEXT)`,
		"CREATE INDEX messages_chat_idx ON messages (chat_id, timestamp)",
	}

	latestVersion = len(dbQueries)
)

type DbConnection struct {
	*sql.DB
}

func Open(driverName, dataSourceName string) (DbConnection, error) {
	db, err := sql.Open(driverName, dataSourceName)
	return DbConnection{db}, err
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
	tx, err := db.Begin()
	for _, sqlStmt := range dbQueries[version:] {
		_, err := tx.Exec(sqlStmt)
		log.Print(sqlStmt)
		if err != nil {
			return err
		}
	}
	err = tx.Commit()
	return err
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
