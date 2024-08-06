package router

import (
	handlers "banner-service/internal/http/handlers"

	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {
	var r = gin.Default()
	r.POST("/banner", handlers.Auth(handlers.AddBannerHandler))
	r.GET("/user_banner", handlers.Auth(handlers.UserBannerHandler))
	r.GET("/banner", handlers.Auth(handlers.UserBannersHanlder))
	r.PATCH("/banner/:id", handlers.Auth(handlers.EditBannerHandler))
	r.DELETE("/banner/:id", handlers.Auth(handlers.EraseBannerHandler))
	return r
}
