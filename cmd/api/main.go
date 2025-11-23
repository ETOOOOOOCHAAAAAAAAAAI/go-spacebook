package main

import (
	"log"

	"SpaceBookProject/internal/auth"
	"SpaceBookProject/internal/config"
	"SpaceBookProject/internal/db"
	"SpaceBookProject/internal/handlers"
	"SpaceBookProject/internal/repository"
	"SpaceBookProject/internal/services"
	"SpaceBookProject/middleware"

	"github.com/gin-gonic/gin"
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("load config: %v", err)
	}
	dbConn, err := db.InitDB(&cfg.Database)
	if err != nil {
		log.Fatalf("init db: %v", err)
	}
	jwtManager := auth.NewJWTManager(cfg.JWT.SecretKey)
	userRepo := repository.NewUserRepository(dbConn)
	spaceRepo := repository.NewSpaceRepository(dbConn)
	authService := services.NewAuthService(userRepo, jwtManager)
	spaceService := services.NewSpaceService(spaceRepo)
	authHandler := handlers.NewAuthHandler(authService)
	spaceHandler := handlers.NewSpaceHandler(spaceService)
	if cfg.Server.Mode != "" {
		gin.SetMode(cfg.Server.Mode)
	}
	r := gin.Default()
	r.Use(middleware.CORSMiddleware())

	api := r.Group(cfg.API.Prefix)
	v1 := api.Group("/" + cfg.API.Version)
	authGroup := v1.Group("/auth")
	{
		authGroup.POST("/register", authHandler.Register)
		authGroup.POST("/login", authHandler.Login)
		authGroup.POST("/refresh", authHandler.RefreshToken)
		authGroup.POST("/logout", authHandler.Logout)
	}
	protected := v1.Group("")
	protected.Use(middleware.AuthMiddleware(jwtManager))
	{
		protected.GET("/auth/me", authHandler.GetMe)
		protected.GET("/spaces", spaceHandler.ListSpaces)
		owner := protected.Group("/owner")
		owner.Use(middleware.OwnerOnlyMiddleware())
		{
			owner.POST("/spaces", spaceHandler.CreateSpace)
			owner.GET("/spaces", spaceHandler.ListMySpaces)
		}
	}

	if err := r.Run(":" + cfg.Server.Port); err != nil {
		log.Fatalf("server run: %v", err)
	}
}
