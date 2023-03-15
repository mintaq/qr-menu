package main

import (
	"os"

	"gitlab.xipat.com/omega-team3/qr-menu-backend/pkg/configs"
	"gitlab.xipat.com/omega-team3/qr-menu-backend/pkg/middleware"
	"gitlab.xipat.com/omega-team3/qr-menu-backend/pkg/routes"
	"gitlab.xipat.com/omega-team3/qr-menu-backend/pkg/utils"
	"gitlab.xipat.com/omega-team3/qr-menu-backend/pkg/worker"
	"gitlab.xipat.com/omega-team3/qr-menu-backend/platform/cache"
	"gitlab.xipat.com/omega-team3/qr-menu-backend/platform/database"

	"github.com/gofiber/fiber/v2"

	_ "gitlab.xipat.com/omega-team3/qr-menu-backend/docs" // load API Docs files (Swagger)

	_ "github.com/joho/godotenv/autoload" // load .env file automatically
)

// @title API
// @version 1.0
// @description This is an auto-generated API Docs.
// @termsOfService http://swagger.io/terms/
// @contact.name API Support
// @contact.email your@mail.com
// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html
// @BasePath /api
// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization
func main() {
	// Define Fiber config.
	config := configs.FiberConfig()
	database.MysqlGormConnection()
	cache.RedisConnection()

	// Define a new Fiber app with config.
	app := fiber.New(config)

	// Middlewares.
	middleware.FiberMiddleware(app) // Register Fiber's middleware for app.

	// Routes.
	routes.SwaggerRoute(app)  // Register a route for API Docs (Swagger).
	routes.PublicRoutes(app)  // Register a public routes for app.
	routes.PrivateRoutes(app) // Register a private routes for app.
	routes.NotFoundRoute(app) // Register route for 404 Error.

	// Redis client.
	worker.CreateRedisClient()
	defer worker.AsynqClient.Close()
	go worker.StartRedisServer()

	// Http client.
	utils.CreateHttpClient()

	// Websocket server
	go utils.StartWebsocketServer()

	// Start server (with or without graceful shutdown).
	if os.Getenv("STAGE_STATUS") == "dev" {
		utils.StartServer(app)
	} else {
		utils.StartServerWithGracefulShutdown(app)
	}
}
