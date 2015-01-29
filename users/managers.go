package users

import (
	"log"
	"database/sql"
	"github.com/linkinpark342/gchat/gchatdb"
)

type UserManager struct {
	db *gchatdb.DbConnection
}

func NewManager(db *gchatdb.DbConnection) *UserManager {
	mgr := UserManager{db}
	return &mgr
}

func (u *UserManager) GetById(id int64) (*User, error) {
	db := u.db
	sqlStmt := "SELECT name FROM users WHERE id = ?"
	var name string
	err := db.QueryRow(sqlStmt, id).Scan(&name)
	switch {
	case err == sql.ErrNoRows:
		log.Printf("No user with that ID")
		return nil, err
	case err != nil:
		log.Fatal(err)
		return nil, err
	default:
		return &User{id: id, name: name}, nil
	}
}

func (u *UserManager) Create(name string) (*User, error) {
	db := u.db
	sqlStmt := "INSERT INTO users(name) VALUES (?)"
	result, err := db.Exec(sqlStmt, name)
	if err != nil {
		return nil, err
	}
	id, err := result.LastInsertId()
	if err != nil {
		return nil, err
	}
	return &User{id: id, name: name}, nil
}
