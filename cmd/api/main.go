package main

import (
	"log"

	"SpaceBookProject/internal/auth"
	"SpaceBookProject/internal/config"
	"SpaceBookProject/internal/db"
	"SpaceBookProject/internal/domain"
	"SpaceBookProject/internal/handlers"
	"SpaceBookProject/internal/repository"
	"SpaceBookProject/internal/services"
	"SpaceBookProject/middleware"

	"github.com/gin-gonic/gin"
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatal(err)
	}

	database, err := db.InitDB(&cfg.Database)
	if err != nil {
		log.Fatal(err)
	}

	jwtManager := auth.NewJWTManager(cfg.JWT.SecretKey)

	userRepo := repository.NewUserRepository(database)
	bookingRepo := repository.NewBookingRepository(database)
	spaceRepo := repository.NewSpaceRepository(database)

	authService := services.NewAuthService(userRepo, jwtManager)
	bookingService := services.NewBookingService(bookingRepo, spaceRepo)
	spaceService := services.NewSpaceService(spaceRepo)

	authHandler := handlers.NewAuthHandler(authService)
	bookingHandler := handlers.NewBookingHandler(bookingService)
	spaceHandler := handlers.NewSpaceHandler(spaceService)

	gin.SetMode(cfg.Server.Mode)
	r := gin.New()
	r.Use(gin.Logger(), gin.Recovery(), middleware.CORSMiddleware())

	api := r.Group(cfg.API.Prefix + "/" + cfg.API.Version)

	authGroup := api.Group("/auth")
	{
		authGroup.POST("/register", authHandler.Register)
		authGroup.POST("/login", authHandler.Login)
		authGroup.POST("/refresh", authHandler.RefreshToken)
		authGroup.GET("/me", middleware.AuthMiddleware(jwtManager), authHandler.GetMe)
	}

	spacesGroup := api.Group("/spaces")
	{
		spacesGroup.GET("", spaceHandler.ListSpaces)
	}
	ownerSpaces := api.Group("/spaces", middleware.AuthMiddleware(jwtManager), middleware.OwnerOnlyMiddleware())
	{
		ownerSpaces.POST("", spaceHandler.CreateSpace)
	}

	bookingsGroup := api.Group("/bookings", middleware.AuthMiddleware(jwtManager))
	{
		bookingsGroup.POST("", middleware.RoleMiddleware(domain.RoleTenant), bookingHandler.CreateBooking)
		bookingsGroup.GET("/my", middleware.RoleMiddleware(domain.RoleTenant), bookingHandler.MyBookings)
		bookingsGroup.PATCH("/:id/cancel", middleware.RoleMiddleware(domain.RoleTenant), bookingHandler.CancelBooking)
	}

	ownerBookings := api.Group("/owner/bookings",
		middleware.AuthMiddleware(jwtManager),
		middleware.OwnerOnlyMiddleware(),
	)
	{
		ownerBookings.GET("", bookingHandler.OwnerBookings)
		ownerBookings.PATCH("/:id/approve", bookingHandler.ApproveBooking)
		ownerBookings.PATCH("/:id/reject", bookingHandler.RejectBooking)
	}

	if err := r.Run(":" + cfg.Server.Port); err != nil {
		log.Fatal(err)
	}
}
