package routes

import (
	"fmt"
	"html"
	"log"
	"os"

	"net/http"

	"github.com/gofiber/fiber/v2"
	socketio "github.com/googollee/go-socket.io"
	"gitlab.xipat.com/omega-team3/qr-menu-backend/app/controllers"
)

// PublicRoutes func for describe group of public routes.
func PublicRoutes(a *fiber.App) {
	// Static routes
	a.Static(os.Getenv("STATIC_PUBLIC_PATH"), "./static/public")

	// Create routes group.
	api := a.Group("/api")
	v1 := api.Group("/v1")

	// Routes for Book:
	v1.Get("/books", controllers.GetBooks)   // get list of all books
	v1.Get("/book/:id", controllers.GetBook) // get one book by ID

	// Routes for User:
	v1.Post("/users/sign/up", controllers.UserSignUp) // register a new user
	v1.Post("/users/sign/in", controllers.UserSignIn) // auth, return Access & Refresh tokens
	v1.Get("/oauth/google/login", controllers.GoogleLogin)
	v1.Get("/oauth/google/callback", controllers.GoogleCallback)
	v1.Post("/users/reset-password", controllers.ResetPassword)
	v1.Post("/users/social/callback", controllers.SocialCallback)

	// Routes for Sapo:
	v1.Get("/sapo/get-token", controllers.GetSapoAccessToken)

	// Routes for Cart:
	v1.Get("/carts", controllers.GetCart)
	v1.Post("/carts/add", controllers.AddItemToCart)
	v1.Post("/carts/update", controllers.UpdateCart)

	// Routes for Order:
	v1.Post("/orders", controllers.CreateOrder)
	v1.Post("/orders/checkout", controllers.CheckoutOrder)
	v1.Post("/orders/pay", controllers.PayOrder)

	// Routes for WS:
	server := socketio.NewServer(nil)

	server.OnConnect("/", func(s socketio.Conn) error {
		s.SetContext("")
		fmt.Println("connected:", s.ID())
		return nil
	})

	server.OnEvent("/", "notice", func(s socketio.Conn, msg string) {
		fmt.Println("notice:", msg)
		s.Emit("reply", "have "+msg)
	})

	server.OnEvent("/chat", "msg", func(s socketio.Conn, msg string) string {
		s.SetContext(msg)
		return "recv " + msg
	})

	server.OnEvent("/", "bye", func(s socketio.Conn) string {
		last := s.Context().(string)
		s.Emit("bye", last)
		s.Close()
		return last
	})

	server.OnError("/", func(s socketio.Conn, e error) {
		fmt.Println("meet error:", e)
	})

	server.OnDisconnect("/", func(s socketio.Conn, reason string) {
		fmt.Println("closed", reason)
	})

	go server.Serve()
	defer server.Close()

	server1 := http.NewServeMux()
	server1.Handle("/socket.io/", server)
	server1.HandleFunc("/test", func(w http.ResponseWriter, r *http.Request) {
		log.Println("alo")
		fmt.Fprintf(w, "Hello, %q", html.EscapeString(r.URL.Path))
	})
	go func() {
		log.Println("Server started on: http://localhost:9001")
		log.Fatalln(http.ListenAndServe("localhost:9001", server1))
	}()
}
