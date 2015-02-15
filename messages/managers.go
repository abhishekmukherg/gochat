package messages

import (
	"time"

	"github.com/linkinpark342/gochat/gchatdb"
	"github.com/linkinpark342/gochat/users"
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

func (c *ChatMgr) GetChat(title string) (*Chat, error) {
	sqlStmt := "SELECT id FROM chats WHERE title = ?"
	var id int64
	err := c.db.QueryRow(sqlStmt, title).Scan(&id)
	if err != nil {
		return nil, err
	}
	sqlStmt = "SELECT user_id FROM chat_users WHERE chat_id = ?"
	rows, err := c.db.Query(sqlStmt, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	users := make([]users.LiteUser, 0)
	for rows.Next() {
		var user_id int64
		err = rows.Scan(&id)
		if err != nil {
			return nil, err
		}
		users = append(users, users.NewLiteUser(user_id))
	}
	return &Chat{Id: id, Title: title, Participants: users}, nil
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
