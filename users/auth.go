package users

import (
	protobuf "github.com/golang/protobuf/proto"
	"github.com/linkinpark342/gochat/proto"
	"golang.org/x/crypto/bcrypt"
	"log"
)

func (u *userModelManager) Authenticate(username string, password []byte) (*User, error) {
	user, err := u.GetByUsername(username)
	if err != nil {
		log.Printf("Failed to get user: %q", err)
		return nil, err
	}
	err = bcrypt.CompareHashAndPassword(user.hashedPassword, password)
	if err != nil {
		log.Printf("Failed to log user in: %q", err)
		return nil, err
	}
	return user, nil
}

func (u *userManager) AuthenticateToken(token string) (user LiteUser) {
	data, err := u.scs.Parse(token)
	if err != nil {
		return nil
	}
	var cookie proto.Cookie
	err = protobuf.Unmarshal(data, &cookie)
	if err != nil {
		log.Fatal("Could not unmarshal a valid cookie:", data)
	}
	return &liteUser{id: *cookie.Id}
}

// Returns a cookie that can validate the user in the future
func (u *tokenManager) GetAuthToken(user *User) string {
	id := user.Id()
	cookie := proto.Cookie{
		Id:          &id,
		AuthVersion: &user.passwordVersion,
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
