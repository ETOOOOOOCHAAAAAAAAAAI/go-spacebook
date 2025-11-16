package main

import (
	"SpaceBookProject/middleware"
	"fmt"
	"log"

	"SpaceBookProject/internal/auth"
	"SpaceBookProject/internal/config"
	"SpaceBookProject/internal/db"
	"SpaceBookProject/internal/handlers"
	"SpaceBookProject/internal/repository"
	"SpaceBookProject/internal/services"

	"github.com/gin-gonic/gin"
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatal("Failed to load config:", err)
	}

	database, err := db.InitDB(&cfg.Database)
	if err != nil {
		log.Fatal("Failed to initialize database:", err)
	}
	defer database.Close()

	userRepo := repository.NewUserRepository(database)

	jwtManager := auth.NewJWTManager(cfg.JWT.SecretKey)

	authService := services.NewAuthService(userRepo, jwtManager)

	authHandler := handlers.NewAuthHandler(authService)

	gin.SetMode(cfg.Server.Mode)
	router := gin.Default()

	router.Use(middleware.CORSMiddleware())

	api := router.Group(cfg.API.Prefix + "/" + cfg.API.Version)
	{
		authRoutes := api.Group("/auth")
		{
			authRoutes.POST("/register", authHandler.Register)
			authRoutes.POST("/login", authHandler.Login)
			authRoutes.POST("/refresh", authHandler.RefreshToken)
		}

		protected := api.Group("")
		protected.Use(middleware.AuthMiddleware(jwtManager))
		{
			protected.GET("/auth/me", authHandler.GetMe)
			protected.POST("/auth/logout", authHandler.Logout)
		}
	}

	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status":  "ok",
			"service": "SpaceBook API",
			"version": cfg.API.Version,
		})
	})

	addr := fmt.Sprintf(":%s", cfg.Server.Port)
	fmt.Printf("Server starting on %s\n", addr)
	fmt.Printf("API endpoints available at %s/%s\n", cfg.API.Prefix, cfg.API.Version)

	if err := router.Run(addr); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
