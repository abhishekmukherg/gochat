package gchatmain

import (
	"fmt"
	"github.com/linkinpark342/gchat/users"
	"github.com/linkinpark342/gchat/gchatdb"
	_ "github.com/mattn/go-sqlite3"
	"log"
)

func Main() {
	db, err := gchatdb.Open("sqlite3", "/tmp/gchat.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	err = db.Upgrade()
	if err != nil {
		log.Fatalf("Failed to upgrade db: %q\n", err)
	}

	userMgr := users.NewManager(db)
	user, err := userMgr.GetById(0)
	fmt.Printf("Magic %v %q\n", user, err)
}
