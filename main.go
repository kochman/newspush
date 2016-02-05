package main

import (
	"golang.org/x/net/websocket"
	"net/http"
)

var cm *connectionManager

func EchoServer(ws *websocket.Conn) {
	cm.register <- ws
	for {
		buf := make([]byte, 256)
		ws.Read(buf)
		cm.broadcast <- buf
	}
}

func StaticServer(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "static")
}

type connectionManager struct {
	connections map[*websocket.Conn]struct{} // empty struct uses zero bytes
	register    chan *websocket.Conn
	broadcast   chan []byte
}

func newConnectionManager() *connectionManager {
	return &connectionManager{
		connections: make(map[*websocket.Conn]struct{}),
		register:    make(chan *websocket.Conn),
		broadcast:   make(chan []byte),
	}
}

func (cm *connectionManager) run() {
	for {
		select {
		case ws := <-cm.register:
			cm.connections[ws] = struct{}{}
		case m := <-cm.broadcast:
			for ws := range cm.connections {
				if _, err := ws.Write(m); err != nil {
					delete(cm.connections, ws)
				}
			}
		}
	}
}

func main() {
	http.HandleFunc("/", StaticServer)
	http.Handle("/echo", websocket.Handler(EchoServer))

	cm = newConnectionManager()
	go cm.run()

	err := http.ListenAndServe(":12345", nil)
	if err != nil {
		panic("ListenAndServe: " + err.Error())
	}
}
