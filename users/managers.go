package users

import (
	"database/sql"
	"errors"
	"github.com/linkinpark342/gchat/gchatdb"
	"github.com/linkinpark342/goscs"
	"golang.org/x/crypto/bcrypt"
	"log"
)

const (
	BCRYPT_COST = 10
)

var (
	ErrMissingField = errors.New("Required field not specified")
)

type UserModelManager interface {
	// Returns a User object from an ID
	GetById(id int64) (*User, error)

	// Returns a User object for a Username
	GetByUsername(username string) (*User, error)

	// Creates a new user and stores in the DB
	Create(name string, password []byte) (*User, error)

	// Returns a user iff the name and password match
	Authenticate(name string, password []byte) (*User, error)
}

// Manages creating and authenticating of tokens, which may be used
// to act as a user's credentials
type TokenManager interface {
	// Returns an auth token for a user
	GetAuthToken(user *User) string

	// Returns a user for an auth token
	AuthenticateToken(token string) (user *User)
}

type UserManager interface {
	TokenManager
	UserModelManager
}

type userModelManager struct {
	db *gchatdb.DbConnection
}

type tokenManager struct {
	scs *goscs.ScsMgr
}

type userManager struct {
	userModelManager
	tokenManager
}

func NewManager(db *gchatdb.DbConnection, scs *goscs.ScsMgr) UserManager {
	umm := userModelManager{db}
	tm := tokenManager{scs}
	mgr := userManager{umm, tm}
	return &mgr
}

func (u *userModelManager) GetById(id int64) (*User, error) {
	sqlStmt := "SELECT id, name, password, passwordVersion FROM users WHERE id = ?"
	return u.getByQuery(sqlStmt, id)
}

func (u *userModelManager) GetByUsername(username string) (*User, error) {
	sqlStmt := "SELECT id, name, password, passwordVersion FROM users WHERE name = ?"
	return u.getByQuery(sqlStmt, username)
}

func (u *userModelManager) getByQuery(sqlStmt string, args ...interface{}) (*User, error) {
	db := u.db
	var id int64
	var name string
	var password []byte
	var passwordVersion int32
	err := db.QueryRow(sqlStmt, args...).Scan(&id, &name, &password, &passwordVersion)
	switch {
	case err == sql.ErrNoRows:
		log.Printf("No user with that ID")
		return nil, nil
	case err != nil:
		log.Fatal(err)
		return nil, err
	default:
		return &User{
			Id:              id,
			Name:            name,
			hashedPassword:  password,
			passwordVersion: passwordVersion}, nil
	}
}

func passwordStrongEnough(password []byte) bool {
	return len(password) >= 8
}

func (u *userModelManager) Create(name string, password []byte) (*User, error) {
	if len(name) == 0 || !passwordStrongEnough(password) {
		return nil, ErrMissingField
	}
	db := u.db
	hashedPassword, err := bcrypt.GenerateFromPassword(password, BCRYPT_COST)
	if err != nil {
		log.Fatalf("Failed to bcrypt password! %q", err)
		return nil, err
	}
	sqlStmt := "INSERT INTO users(name, password, passwordVersion) VALUES (?, ?, 1)"
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
