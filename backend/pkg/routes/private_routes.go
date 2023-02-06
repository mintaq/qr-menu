package routes

import (
	"github.com/gofiber/fiber/v2"
	"gitlab.xipat.com/omega-team3/qr-menu-backend/app/controllers"
	"gitlab.xipat.com/omega-team3/qr-menu-backend/pkg/middleware"
)

// PrivateRoutes func for describe group of private routes.
func PrivateRoutes(a *fiber.App) {
	// Create routes group.
	api := a.Group("/api")
	v1 := api.Group("/v1", middleware.JWTProtected())

	// Routes for Token:
	v1.Post("/token/renew", controllers.RenewTokens) // renew Access & Refresh tokens

	// Routes for Book:
	v1.Post("/book", controllers.CreateBook)   // create a new book
	v1.Delete("/book", controllers.DeleteBook) // delete one book by ID
	v1.Put("/book", controllers.UpdateBook)    // update one book by ID

	// Routes for User:
	v1.Post("/user/sign/out", controllers.UserSignOut) // de-authorization user
	v1.Post("/user/create-password", controllers.CreateNewPassword)

	// Routes for Kiotviet:
	v1.Post("/kiotviet/create-user", controllers.CreateKiotvietUser)

	// Routes for Business:
	v1.Post("/business/create", controllers.CreateBusiness)
}
