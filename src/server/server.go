package server

import (
	"os"

	"github.com/Yukio0315/reservation-backend/src/controller"
	"github.com/Yukio0315/reservation-backend/src/middleware"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

// Init runs gin router
func Init() {
	r := new()
	r = router(r)
	r.Run()
}

func new() *gin.Engine {
	port := os.Getenv("PORT")
	r := gin.New()

	config := cors.DefaultConfig()
	config.AllowOrigins = []string{"http://localhost:3000"}
	config.AllowHeaders = []string{"Authorization"}

	r.Use(cors.New(config))
	r.Use(gin.Logger())
	r.Use(gin.Recovery())

	if port == "" {
		port = "8000"
	}
	return r
}

func router(r *gin.Engine) *gin.Engine {

	authMiddleware := middleware.AuthMiddleware()
	userCtrl := controller.UserController{}
	eventCtrl := controller.EventController{}
	reservationCtrl := controller.ReservationController{}

	a := r.Group("/v1")
	a.POST("/users", userCtrl.Create)
	a.POST("/sign-in", authMiddleware.LoginHandler)
	a.POST("/reset-password", userCtrl.ReserveResetPassword)
	a.GET("/reset-password/:uuid", userCtrl.CheckUUID)
	a.PATCH("/reset-password/:uuid", userCtrl.PasswordReset)

	auth := a.Group("")
	auth.Use(authMiddleware.MiddlewareFunc())
	{
		e := auth.Group("/events/:id")
		{
			e.GET("", eventCtrl.Show)
		}

		r := auth.Group("/reservations/:id")
		{
			r.POST("", reservationCtrl.Add)
			r.DELETE("", reservationCtrl.Cancel)
		}

		u := auth.Group("/users/:id")
		{
			u.GET("", userCtrl.Show)
			u.DELETE("", userCtrl.Delete)
			u.PATCH("/password", userCtrl.PasswordChange)
			u.PATCH("/username", userCtrl.UserNameChange)
			u.PATCH("/email", userCtrl.EmailChange)
		}

		admin := auth.Group("/admin/:id")
		admin.Use(middleware.AdminMiddleware())
		{
			admin.DELETE("/events", eventCtrl.Delete)
		}
	}

	return r
}
