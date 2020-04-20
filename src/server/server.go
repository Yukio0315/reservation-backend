package server

import (
	"os"

	"github.com/Yukio0315/reservation-backend/src/controller"
	"github.com/Yukio0315/reservation-backend/src/middleware"
	"github.com/gin-gonic/gin"
)

// Init runs gin router
func Init() {
	r := router()
	r.Run()
}

func router() *gin.Engine {
	port := os.Getenv("PORT")
	r := gin.New()
	r.Use(gin.Logger())
	r.Use(gin.Recovery())

	if port == "" {
		port = "8000"
	}

	authMiddleware := middleware.AuthMiddleware()

	userCtrl := controller.UserController{}
	eventCtrl := controller.EventController{}
	reservationCtrl := controller.ReservationController{}

	a := r.Group("/api")
	a.POST("/signin", authMiddleware.LoginHandler)
	a.POST("/login", authMiddleware.LoginHandler)
	a.PATCH("/reset-password", userCtrl.PasswordReset)

	auth := a.Group("")
	auth.Use(authMiddleware.MiddlewareFunc())
	{
		e := auth.Group("/event/:id")
		{
			e.GET("", eventCtrl.Show)
		}
		r := auth.Group("/reservation/:id")
		{
			r.POST("/", reservationCtrl.Add)
			r.DELETE("/", reservationCtrl.Cancel)
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
			admin.DELETE("/event", eventCtrl.Delete)
		}
	}

	return r
}
