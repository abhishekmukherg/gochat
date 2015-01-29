package users

import (
	"fmt"
)

type User struct {
	id int64
	name string
}

func (u *User) GetId() int64 {
	return u.id
}

func (u *User) GetName() string {
	return u.name
}

func (u *User) String() string {
	return fmt.Sprintf("User{%v, %v}", u.id, u.name)
}
