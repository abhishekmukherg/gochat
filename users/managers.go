package users

import (
	"database/sql"
	"errors"
	"github.com/linkinpark342/gchat/gchatdb"
	"golang.org/x/crypto/bcrypt"
	"log"
)

const (
	BCRYPT_COST = 10
)

var (
	ErrMissingField = errors.New("Required field not specified")
)

type UserManager struct {
	db *gchatdb.DbConnection
}

func NewManager(db *gchatdb.DbConnection) *UserManager {
	mgr := UserManager{db}
	return &mgr
}

func (u *UserManager) GetById(id int64) (*User, error) {
	sqlStmt := "SELECT id, name, password FROM users WHERE id = ?"
	return u.getByQuery(sqlStmt, id)
}

func (u *UserManager) GetByUsername(username string) (*User, error) {
	sqlStmt := "SELECT id, name, password FROM users WHERE name = ?"
	return u.getByQuery(sqlStmt, username)
}

func (u *UserManager) getByQuery(sqlStmt string, args ...interface{}) (*User, error) {
	db := u.db
	var id int64
	var name string
	var password []byte
	err := db.QueryRow(sqlStmt, args...).Scan(&id, &name, &password)
	switch {
	case err == sql.ErrNoRows:
		log.Printf("No user with that ID")
		return nil, nil
	case err != nil:
		log.Fatal(err)
		return nil, err
	default:
		return &User{Id: id, Name: name, hashedPassword: password}, nil
	}
}

func passwordStrongEnough(password []byte) bool {
	return len(password) >= 8
}

func (u *UserManager) Create(name string, password []byte) (*User, error) {
	if len(name) == 0 || !passwordStrongEnough(password) {
		return nil, ErrMissingField
	}
	db := u.db
	hashedPassword, err := bcrypt.GenerateFromPassword(password, BCRYPT_COST)
	if err != nil {
		log.Fatalf("Failed to bcrypt password! %q", err)
		return nil, err
	}
	sqlStmt := "INSERT INTO users(name, password) VALUES (?, ?)"
	result, err := db.Exec(sqlStmt, name, hashedPassword)
	if err != nil {
		return nil, err
	}
	id, err := result.LastInsertId()
	if err != nil {
		return nil, err
	}
	return &User{Id: id, Name: name}, nil
}
