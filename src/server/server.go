package server

import (
	user "github.com/Yukio0315/reservation-backend/src/controller"
	"github.com/gin-gonic/gin"
)

func Init() {
	r := router()
	r.Run()
}

func router() *gin.Engine {
	r := gin.Default()

	u := r.Group("/users")
	{
		ctrl := user.Controller{}
		u.GET("", ctrl.Index)
	}

	return r
}
