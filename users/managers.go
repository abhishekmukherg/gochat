package users

import (
	"strconv"
)

type UserManager uint64

func NewManager() *UserManager {
	mgr := UserManager(0)
	return &mgr
}

func (*UserManager) GetById(id uint64) *User {
	return &User{id: id, name: strconv.FormatUint(id, 10)}
}
