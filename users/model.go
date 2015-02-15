package users

import (
	"fmt"
)

type LiteUser interface {
	Id() int64
}

type liteUser struct {
	id int64
}

func (l liteUser) Id() int64 {
	return l.id
}

func NewLiteUser(id int64) LiteUser {
	return liteUser{id}
}

type User struct {
	liteUser
	Name            string
	hashedPassword  []byte
	passwordVersion int32
}

func (u *User) String() string {
	return fmt.Sprintf("User{%v, %v}", u.Id, u.Name)
}
