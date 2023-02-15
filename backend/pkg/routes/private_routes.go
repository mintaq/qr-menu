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

	// Routes for Sapo:
	v1.Get("/sapo/get-auth-url", controllers.GetSapoAuthURL)
	v1.Post("/sapo/products/sync", controllers.SyncSapoProducts)
	v1.Post("/sapo/custom-collections/sync", controllers.SyncSapoCustomCollections)
	v1.Post("/sapo/smart-collections/sync", controllers.SyncSapoSmartCollections)

	// Routes for Store:
	v1.Post("/store/create", controllers.CreateStore)
	v1.Get("/store", controllers.GetStore)

	// Routes for Theme:
	v1.Post("/theme/create", controllers.CreateTheme)
	v1.Get("/themes", controllers.GetThemes)

	// Routes for QR Code:
	v1.Post("qr-code/create", controllers.CreateQrCode)

	// Routes for App:
	v1.Post("/app/create", controllers.CreateApp)

	// Routes for Product:
	v1.Post("/product/create", controllers.CreateProduct)
	v1.Get("/products", controllers.GetProducts)

	// Routes for Menu:
	v1.Post("/menu/create", controllers.CreateMenu)
	v1.Get("/menus", controllers.GetMenus)
}
