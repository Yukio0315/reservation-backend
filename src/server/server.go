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
	// reservationCtrl := controller.ReservationController{}

	r.POST("/signin", authMiddleware.LoginHandler)
	r.POST("/login", authMiddleware.LoginHandler)

	auth := r.Group("/auth")
	auth.Use(authMiddleware.MiddlewareFunc())
	{
		// r := auth.Group("/reservation")
		// {
		// 	r.GET("/", reservationCtrl.Show)
		// 	r.GET("/", reservationCtrl.Add)
		// }

		u := auth.Group("/users/:id")
		{
			u.GET("", userCtrl.Show)
			// u.PATCH("/password", userCtrl.UpdatePassword)
			// u.PATCH("/user-name", userCtrl.UpdateUserName)
			// u.PATCH("/email", userCtrl.UpdateEmail)
			// u.DELETE("/cancel", userCtrl.Cancel)
			// u.DELETE("/delete", userCtrl.Delete)
			// u.DELETE("/logout", userCtrl.Logout)
		}
	}

	return r
}
