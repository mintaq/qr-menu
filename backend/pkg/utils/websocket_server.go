package utils

import (
	"fmt"
	"html"
	"log"
	"net/http"

	socketio "github.com/googollee/go-socket.io"
)

func StartWebsocketServer() {
	// Routes for WS:
	server := socketio.NewServer(nil)

	server.OnConnect("/", func(s socketio.Conn) error {
		s.SetContext("")
		room := fmt.Sprintf("some:room:%v", s.ID())
		s.Join(room)
		fmt.Println("joined:", room)
		return nil
	})

	server.OnEvent("/notice", "notice", func(s socketio.Conn, msg string) {
		fmt.Println("notice:", msg)
		s.Emit("reply", "have "+msg)
	})

	server.OnEvent("/chat", "msg", func(s socketio.Conn, msg string) string {
		s.SetContext(msg)
		return "recv " + msg
	})

	server.OnEvent("/bye", "bye", func(s socketio.Conn) string {
		last := s.Context().(string)
		s.Emit("bye", last)
		s.Close()
		return last
	})

	server.OnError("/error", func(s socketio.Conn, e error) {
		fmt.Println("meet error:", e)
	})

	server.OnDisconnect("/dis", func(s socketio.Conn, reason string) {
		fmt.Println("closed", reason)
	})

	go server.Serve()
	defer server.Close()

	serverMux := http.NewServeMux()
	serverMux.Handle("/socket.io/", server)
	serverMux.HandleFunc("/hello", CORS(func(w http.ResponseWriter, r *http.Request) {
		log.Println("alo")
		fmt.Fprintf(w, "Hello, %q", html.EscapeString(r.URL.Path))
	}))
	log.Println("Server started on: http://localhost:5009")
	log.Fatalln(http.ListenAndServe(":5009", serverMux))
}
