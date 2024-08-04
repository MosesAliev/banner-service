package router

import (
	handlers "banner-service/internal/http/handlers"

	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {
	var r = gin.Default()
	r.POST("/banner", handlers.AddBannerHandler)
	r.GET("/user_banner", handlers.UserBannerHandler)
	r.GET("/banner", handlers.UserBannersHanlder)
	r.PATCH("/banner/:id", handlers.EditBannerHandler)
	r.DELETE("/banner/:id", handlers.EraseBannerHandler)
	return r
}
