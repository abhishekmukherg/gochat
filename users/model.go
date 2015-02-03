package users

import (
	"fmt"
)

type LiteUser struct {
	Id int64
}

type User struct {
	LiteUser
	Name            string
	hashedPassword  []byte
	passwordVersion int32
}

func (u *User) String() string {
	return fmt.Sprintf("User{%v, %v}", u.Id, u.Name)
}
