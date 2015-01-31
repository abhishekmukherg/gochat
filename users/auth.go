package users

import (
	protobuf "github.com/golang/protobuf/proto"
	"github.com/linkinpark342/gchat/proto"
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

func (u *tokenManager) AuthenticateToken(token string) (user *User) {
	data, err := u.scs.Parse(token)
	if err != nil {
		return nil
	}
	var cookie proto.Cookie
	err = protobuf.Unmarshal(data, &cookie)
	if err != nil {
		log.Fatal("Could not unmarshal a valid cookie:", data)
	}
	user = new(User)
	user.Id = *cookie.Id
	return user
}

// Returns a cookie that can validate the user in the future
func (u *tokenManager) GetAuthToken(user *User) string {
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
