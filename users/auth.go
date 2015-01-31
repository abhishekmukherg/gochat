package users

import (
	"golang.org/x/crypto/bcrypt"
	"log"
)

func (u *UserManager) Authenticate(username string, password []byte) *User {
	user, err := u.GetByUsername(username)
	if err != nil {
		log.Printf("Failed to get user: %q", err)
		return nil
	}
	err = bcrypt.CompareHashAndPassword(user.hashedPassword, password)
	if err != nil {
		log.Printf("Failed to log user in: %q", err)
		return nil
	}
	return user
}
