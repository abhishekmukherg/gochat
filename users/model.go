package users

import (
	"fmt"
)

type User struct {
	Id             int64
	Name           string
	hashedPassword []byte
}

func (u *User) String() string {
	return fmt.Sprintf("User{%v, %v}", u.Id, u.Name)
}
