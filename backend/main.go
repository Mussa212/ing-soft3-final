package main

import (
	"log"
	"net/url"

	"github.com/gin-gonic/gin"

	"vesuvio/internal/client"
	"vesuvio/internal/config"
	"vesuvio/internal/controller"
	"vesuvio/internal/middleware"
	"vesuvio/internal/service"
)

var (
	openDB    = client.NewDB
	startHTTP = func(r *gin.Engine, port string) error { return r.Run(":" + port) }
)

func main() {
	cfg := config.Load()

	log.Printf("Using DATABASE_DSN: %s", redactDSN(cfg.DatabaseDSN))
	db, err := openDB(cfg.DatabaseDSN)
	if err != nil {
		log.Fatalf("failed to connect database: %v", err)
	}

	userClient := client.NewUserClient(db)
	reservationClient := client.NewReservationClient(db)

	authService := service.NewAuthService(userClient)
	reservationService := service.NewReservationService(reservationClient)

	authController := controller.NewAuthController(authService)
	reservationController := controller.NewReservationController(reservationService)
	adminController := controller.NewAdminController(reservationService)

	r := gin.Default()
	r.Use(middleware.CORSMiddleware())

	r.POST("/auth/register", authController.Register)
	r.POST("/auth/login", authController.Login)

	authRequired := r.Group("/")
	authRequired.Use(middleware.AuthMiddleware(authService))
	{
		authRequired.GET("/my/reservations", reservationController.ListMyReservations)
		authRequired.POST("/reservations", reservationController.CreateReservation)
		authRequired.PATCH("/reservations/:id/cancel", reservationController.CancelReservation)
	}

	adminRequired := r.Group("/admin")
	adminRequired.Use(middleware.AuthMiddleware(authService), middleware.AdminOnly())
	{
		adminRequired.GET("/reservations", adminController.ListReservations)
		adminRequired.PATCH("/reservations/:id/confirm", adminController.ConfirmReservation)
		adminRequired.PATCH("/reservations/:id/cancel", adminController.CancelReservation)
	}

	if err := startHTTP(r, cfg.Port); err != nil {
		log.Fatalf("failed to start server: %v", err)
	}
}

// redactDSN masks the password in the DSN for logging.
func redactDSN(dsn string) string {
	u, err := url.Parse(dsn)
	if err != nil || u.User == nil {
		return dsn
	}
	username := u.User.Username()
	u.User = url.UserPassword(username, "****")
	return u.String()
}
