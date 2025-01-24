package main

import (
	"log"
	"net/http"
	"gorm.io/gorm"
	"github.com/vercel/go-bridge/go/bridge"
	"github.com/gin-gonic/gin"
	
	"nevacarwash.com/main/database"
	"nevacarwash.com/main/handlers"
	"nevacarwash.com/main/middleware"
	"nevacarwash.com/main/repositories"
	"nevacarwash.com/main/services"
)

var db *gorm.DB

func init() {
	database.LoadEnvs()
	database.InitializeDatabaseLayer()

	// Check if tables exist first
	if !database.TablesExist() {
		log.Println("Tables do not exist. Running migrations...")
		if err := database.Migrate(); err != nil {
			log.Fatalf("Failed to migrate database: %v", err)
		}
		log.Println("Migrations completed successfully")
	} else {
		log.Println("Tables already exist. Skipping migrations")
	}
	db = database.GetDB()
}

func createRouter() *gin.Engine {
	router := gin.Default()

	// Load templates and static files
	router.LoadHTMLGlob("templates/*")
	router.Static("/static", "./static")

	// Define routes
	router.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "home.html", nil)
	})

	router.GET("/api/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})

	return router
}

// Handler is the exported function Vercel uses as the entry point
func handler(w http.ResponseWriter, r *http.Request) {
    router := createRouter()
    bridge.ServeHTTP(router, w, r) // Use Vercel's Go bridge to handle requests
}

func main() {
	// Create repository
	vehicleRepo := repositories.NewVehicleRepository(db)

	// Create service
	vehicleService := services.NewVehicleService(vehicleRepo)

	// Create handler
	vehicleHandler := handlers.NewVehicleHandler(vehicleService)

	// setup gin router
	router := gin.Default()
	router.Use(gin.Logger())

	// Load HTML templates
	router.LoadHTMLGlob("templates/*")

	// Auth routes
	auth := router.Group("/")
	{
		auth.GET("", handlers.Home)
		auth.GET("/login", handlers.Login)
		auth.GET("/logout", handlers.Logout)
		auth.GET("/register", handlers.CreateUser)
		auth.POST("/login", handlers.Login)
		auth.POST("/register", handlers.CreateUser)
	}

	// Vehicle routes
	snip := router.Group("/vehicles")
	{
		// Guest routes
		snip.GET("", vehicleHandler.GetVehiclesByStatus)
		snip.GET("/:id", vehicleHandler.GetVehicleByID)

		// Authenticated routes
		snip.GET("/new", middleware.CheckAuth, vehicleHandler.CreateVehicle)
		snip.POST("/new", middleware.CheckAuth, vehicleHandler.CreateVehicle)
		snip.GET("/:id/edit", middleware.CheckAuth, vehicleHandler.UpdateVehicle)
		snip.POST("/:id/edit", middleware.CheckAuth, vehicleHandler.UpdateVehicle)
		snip.POST("/:id/delete", middleware.CheckAuth, vehicleHandler.DeleteVehicle)
		snip.GET("/:id/delete", middleware.CheckAuth, vehicleHandler.DeleteVehicle)
	}

	// start server
	log.Println("starting server on :8080")
	if err := router.Run(":8080"); err != nil {
		log.Fatalf("failed to start server: %v", err)
	}
}
