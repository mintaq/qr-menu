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
	v1.Post("/tokens/renew", controllers.RenewTokens) // renew Access & Refresh tokens

	// Routes for Book:
	v1.Post("/book", controllers.CreateBook)   // create a new book
	v1.Delete("/book", controllers.DeleteBook) // delete one book by ID
	v1.Put("/book", controllers.UpdateBook)    // update one book by ID

	// Routes for User:
	v1.Post("/users/sign/out", controllers.UserSignOut) // de-authorization user
	v1.Post("/users/create-password", controllers.CreateNewPassword)
	v1.Get("/users/profile", controllers.GetUserProfile)

	// Routes for Kiotviet:
	v1.Post("/kiotviet/create-user", controllers.CreateKiotvietUser)
	v1.Post("/kiotviet/products/sync", controllers.SyncKiotvietProducts)
	v1.Post("/kiotviet/collections/sync", controllers.SyncKiotvietCollections)

	// Routes for Sapo:
	v1.Get("/sapo/get-auth-url", controllers.GetSapoAuthURL)
	v1.Post("/sapo/products/sync", controllers.SyncSapoProducts)
	v1.Post("/sapo/custom-collections/sync", controllers.SyncSapoCustomCollections)
	v1.Post("/sapo/smart-collections/sync", controllers.SyncSapoSmartCollections)

	// Routes for Store:
	v1.Post("/stores", controllers.CreateStore)
	v1.Get("/stores", controllers.GetStores)
	v1.Get("/stores/:id", controllers.GetStoreById)
	v1.Get("/stores/subdomain/:subdomain", controllers.GetStoreBySubdomain)
	v1.Put("/stores/:id", controllers.UpdateStore)

	// Routes for Theme:
	v1.Post("/themes", controllers.CreateTheme)
	v1.Get("/themes", controllers.GetThemes)
	v1.Get("/themes/main", controllers.GetMainTheme)
	v1.Put("/themes/:id", controllers.UpdateTheme)

	// Routes for App:
	v1.Post("/apps", controllers.CreateApp)

	// Routes for Product:
	v1.Post("/products", controllers.CreateProduct)
	v1.Get("/products", controllers.GetProducts)
	v1.Get("/products/:product_id", controllers.GetProductByProductId)
	v1.Delete("/products/:product_id", controllers.DeleteProduct)
	v1.Put("/products/:product_id", controllers.UpdateProduct)

	// Routes for Collection:
	v1.Get("/collections", controllers.GetCollections)
	v1.Get("/collections/featured", controllers.GetFeaturedCollection)
	v1.Post("/collections", controllers.CreateCollection)

	// Routes for Menu:
	v1.Post("/menus", controllers.CreateMenu)
	v1.Get("/menus", controllers.GetMenus)
	v1.Get("/menus/main", controllers.GetMainMenu)
	v1.Delete("/menus/:id", controllers.DeleteMenu)
	v1.Put("/menus/:id", controllers.UpdateMenu)

	// Routes for Table:
	v1.Post("/tables", controllers.CreateTable)
	v1.Get("/tables", controllers.GetTables)
	v1.Delete("/tables/:id", controllers.DeleteTable)

	// Route for Collect:
	v1.Post("/collects", controllers.CreateCollect)
}
