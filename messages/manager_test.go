package messages

import (
	"testing"

	"github.com/linkinpark342/gchat/gchatdb"
	"github.com/linkinpark342/gchat/users"
	"github.com/linkinpark342/goscs"
	_ "github.com/mattn/go-sqlite3"
	. "gopkg.in/check.v1"
)

func Test(t *testing.T) { TestingT(t) }

type MySuite struct {
	db gchatdb.DbConnection
	cm ChatMgr
	um users.UserManager
}

var _ = Suite(&MySuite{})

func (s *MySuite) SetUpTest(c *C) {
	var err error
	// set up db
	s.db, err = gchatdb.Open("sqlite3", ":memory:")
	c.Assert(err, IsNil)
	err = s.db.Upgrade()
	c.Assert(err, IsNil)
	scsMgr := goscs.NewMgr([]byte("deadbedwasfed123"))
	// set up chatmgr
	s.cm = NewManager(s.db)
	c.Assert(err, IsNil)
	// set up usermgr
	s.um = users.NewManager(s.db, scsMgr)
	c.Assert(err, IsNil)
}

func (s *MySuite) TearDownTest(c *C) {
	err := s.db.Close()
	c.Assert(err, IsNil)
}

func (s *MySuite) TestCreateChat(c *C) {
	chat, err := s.cm.NewChat("devops")
	c.Assert(err, IsNil)
	c.Assert(chat, NotNil)
}

func (s *MySuite) TestAddParticipant(c *C) {
	chat, _ := s.cm.NewChat("devops")
	user, err := s.um.Create("username", []byte("password"))
	c.Assert(err, IsNil)
	chat, err = s.cm.Subscribe(chat, *user)
	c.Check(err, IsNil)
	c.Check(chat.Participants, HasLen, 1)
	c.Check(chat.Participants, DeepEquals, []users.LiteUser{*user})
}
