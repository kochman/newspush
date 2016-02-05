package main

import (
	"io"
	"net/http"
	"golang.org/x/net/websocket"
)

func EchoServer(ws *websocket.Conn) {
	io.Copy(ws, ws)
}

func StaticServer(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "static")
}

func main() {
	http.HandleFunc("/", StaticServer)
	http.Handle("/echo", websocket.Handler(EchoServer))
	err := http.ListenAndServe(":12345", nil)
	if err != nil {
		panic("ListenAndServe: " + err.Error())
	}
}
