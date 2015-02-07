package messages

import (
	"time"

	"github.com/linkinpark342/gchat/gchatdb"
	"github.com/linkinpark342/gchat/users"
)

type ChatMgr struct {
	db gchatdb.DbConnection
}

func NewManager(db gchatdb.DbConnection) ChatMgr {
	return ChatMgr{db}
}

func (c *ChatMgr) NewChat(title string) (*Chat, error) {
	if len(title) == 0 {
		return nil, gchatdb.ErrMissingField
	}
	sqlStmt := "INSERT INTO chats(title) VALUES (?)"
	result, err := c.db.Exec(sqlStmt, title)
	if err != nil {
		return nil, err
	}
	id, err := result.LastInsertId()
	if err != nil {
		return nil, err
	}
	return &Chat{Id: id, Title: title}, nil
}

func (cm *ChatMgr) Subscribe(c *Chat, users ...users.LiteUser) (*Chat, error) {
	tx, err := cm.db.Begin()
	if err != nil {
		return c, err
	}
	sqlStmt := "INSERT INTO chat_users(chat_id, user_id) VALUES (?, ?)"
	prepared, err := tx.Prepare(sqlStmt)
	if err != nil {
		return c, err
	}
	for _, u := range users {
		_, err = prepared.Exec(c.Id, u.Id())
		if err != nil {
			tx.Rollback()
			return c, err
		}
	}
	err = tx.Commit()
	if err != nil {
		return c, err
	}
	c.Participants = append(c.Participants, users...)
	return c, nil
}

func (cm *ChatMgr) NewMessage(c *Chat, u users.LiteUser, text string) (*Message, error) {
	now := time.Now()
	sqlStmt := "INSERT INTO messages(user_id, chat_id, timestamp, text) VALUES (?, ?, ?, ?)"
	result, err := cm.db.Exec(sqlStmt, u.Id, c.Id, now, text)
	if err != nil {
		return nil, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return nil, err
	}
	return &Message{
		Id:        id,
		UserId:    u.Id(),
		ChatId:    c.Id,
		Timestamp: now,
		Text:      text,
	}, nil
}
