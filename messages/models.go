package messages

import (
	"fmt"
	"github.com/linkinpark342/gochat/users"
	"time"
)

type Chat struct {
	Id           int64
	Title        string
	Participants []users.LiteUser
}

type Message struct {
	Id        int64
	UserId    int64
	ChatId    int64
	Timestamp time.Time
	Text      string
}

func (c Chat) String() string {
	return fmt.Sprintf("Chat{Id: %d, Title: %s}", c.Id, c.Title)
}
