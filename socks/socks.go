package main

import (
	"io"
	"net/http"

	"golang.org/x/net/websocket"
)

type T struct {
	Msg   string
	Count int
}

// receive JSON type T
func main() {
	http.Handle("/", websocket.Handler(func(c *websocket.Conn) {
		io.Copy(c, c)
	}))
	http.ListenAndServe(":12345", nil)
}
