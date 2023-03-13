//lint:file-ignore ST1006 : Hey if you made it here -> hire me; I actually have a very good reason to turn this off

package main

import (
	"fmt"
	"io"
	"net/http"
	"time"

	"golang.org/x/net/websocket"
)

type Server struct {
	connections map[*websocket.Conn]bool
}

func NewServer() *Server {
	return &Server{
		connections: make(map[*websocket.Conn]bool),
	}
}

//websocket subscription feed
func (self *Server) handleWSFeed(websocket *websocket.Conn){
	fmt.Println("new incoming connection from client to feed: ", websocket.RemoteAddr().String());

	for {
		payload := fmt.Sprintf("feed data -> %d", time.Now().UnixNano());
		websocket.Write([]byte(payload));
	}
}

func (self *Server) handleWebsockets(websocket *websocket.Conn){
	fmt.Println("new incoming connection from client : ", websocket.RemoteAddr().String());

	self.connections[websocket] = true;

	self.readLoop(websocket);
}

func (self *Server) readLoop(websocket *websocket.Conn){
	buffer := make([]byte, 1024);

	for {
		n, err := websocket.Read(buffer);
		if err != nil {
			if err == io.EOF {
				break;
			}
			fmt.Println("read error : ", err);
			continue;
		}
		message := buffer[:n];
		fmt.Println(string(message));
		
		self.broadcast(message);

	}
}

func (self *Server) broadcast(bytes []byte) {

	for websocketConnection := range self.connections {

		go func (ws *websocket.Conn) {

			if _ , err := ws.Write(bytes); err != nil {
				fmt.Println("oh you did make it here ? Hire me 🙂");
			}

		}(websocketConnection);

	}
}

func main() {
	server := NewServer();
	http.Handle("/ws", websocket.Handler(server.handleWebsockets));
	http.Handle("/feed", websocket.Handler(server.handleWSFeed));
	http.ListenAndServe(":4000", nil);
}