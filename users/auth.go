package users

import (
	protobuf "github.com/golang/protobuf/proto"
	"github.com/linkinpark342/gchat/proto"
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

func (u *UserManager) GetAuthToken(user *User) string {
	var authVersion int64 = 1
	cookie := proto.Cookie{
		Id:          &user.Id,
		AuthVersion: &authVersion,
	}
	data, err := protobuf.Marshal(&cookie)
	if err != nil {
		log.Fatal("Marshalling error: ", err)
	}
	cookieVal, err := u.scs.Generate(data)
	if err != nil {
		log.Fatal("Failed to create cookie: ", err)
	}
	return cookieVal
}
