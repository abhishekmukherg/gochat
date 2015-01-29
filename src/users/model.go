package users

import (
	"fmt"
)

type User struct {
	id uint64
	name string
}

func (u *User) getId() uint64 {
	return u.id
}

func (u *User) getName() string {
	return u.name
}

func (u *User) String() string {
	return fmt.Sprintf("User{%v, %v}", u.id, u.name)
}
