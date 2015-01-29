package gchatmain

import (
	"github.com/linkinpark342/gchat/gchatdb"
	"github.com/linkinpark342/gchat/users"
	"github.com/linkinpark342/gchat/router"
	_ "github.com/mattn/go-sqlite3"
	"net/http"
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

	handler := router.Create(userMgr)
	http.Handle("/", handler)
	http.ListenAndServe(":8080", nil)
}
