package routes

import (
	"github.com/gofiber/fiber/v2"
	"gitlab.xipat.com/omega-team3/qr-menu-backend/app/controllers"
)

// PublicRoutes func for describe group of public routes.
func PublicRoutes(a *fiber.App) {
	// Create routes group.
	api := a.Group("/api")
	v1 := api.Group("/v1")

	// Routes for Book:
	v1.Get("/books", controllers.GetBooks)   // get list of all books
	v1.Get("/book/:id", controllers.GetBook) // get one book by ID

	// Routes for User:
	v1.Post("/user/sign/up", controllers.UserSignUp) // register a new user
	v1.Post("/user/sign/in", controllers.UserSignIn) // auth, return Access & Refresh tokens
	v1.Get("/oauth/google/login", controllers.GoogleLogin)
	v1.Get("/oauth/google/callback", controllers.GoogleCallback)
	v1.Post("/user/reset-password", controllers.ResetPassword)
}
