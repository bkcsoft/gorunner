package main

import (
	"net/http"

	"github.com/gorilla/websocket"
	. "github.com/jakecoffman/gorunner/service"
)

var nothing = map[string]string{}

func errHelp(msg string) map[string]interface{} {
	return map[string]interface{}{"error": msg}
}

// General

func app(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "web/static/app.html")
}

func favicon(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "web/favicon.ico")
}

func wsHandler(c context, w http.ResponseWriter, r *http.Request) (int, interface{}) {
	// Upgrade the HTTP connection to a websocket
	ws, err := websocket.Upgrade(w, r, nil, 1024, 1024)
	if _, ok := err.(websocket.HandshakeError); ok {
		return http.StatusBadRequest, errHelp("Not a websocket handshake")
	} else if err != nil {
		return http.StatusInternalServerError, errHelp(err.Error())
	}
	conn := NewConnection(ws)
	c.Hub().Register(conn)
	defer c.Hub().Unregister(conn)
	go conn.Writer()
	conn.Reader()
	return http.StatusOK, nothing
}
