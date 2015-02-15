package gchatmain

import (
	"github.com/linkinpark342/gochat/gchatdb"
	"github.com/linkinpark342/gochat/router"
	"github.com/linkinpark342/gochat/users"
	"github.com/linkinpark342/goscs"
	_ "github.com/mattn/go-sqlite3"
	"log"
	"net/http"
)

func Main() {
	db, err := gchatdb.Open("sqlite3", "/tmp/gchat.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	scsMgr := goscs.NewMgr([]byte("deadbedwasfed123"))

	err = db.Upgrade()
	if err != nil {
		log.Fatalf("Failed to upgrade db: %q\n", err)
	}

	userMgr := users.NewManager(db, scsMgr)

	handler := router.Create(userMgr)
	http.Handle("/", handler)
	http.ListenAndServe(":8080", nil)
}
