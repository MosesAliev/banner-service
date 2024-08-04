package handlers

import (
	"banner-service/internal/database"
	"banner-service/internal/models"
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// Получение баннера для пользователя
func UserBannerHandler(c *gin.Context) {
	var tagID, tagOk = c.GetQuery("tag_id")
	if !tagOk {
		log.Println("Отсутствует тег в запросе")
		c.IndentedJSON(http.StatusBadRequest, models.ErrorResponse{Error: "string"})

		return
	}

	var featureID, featureOk = c.GetQuery("feature_id")
	if !featureOk {
		log.Println("Отсутствует фича в запросе")
		c.IndentedJSON(http.StatusBadRequest, models.ErrorResponse{Error: "string"})

		return
	}

	var tagFeatureBanner = models.TagFeatureBanner{}

	var err error
	tagFeatureBanner.TagID, err = strconv.Atoi(tagID)
	if err != nil {
		log.Println(err)
		c.IndentedJSON(http.StatusBadRequest, models.ErrorResponse{Error: "string"})

		return
	}

	tagFeatureBanner.FeatureID, err = strconv.Atoi(featureID)
	if err != nil {
		log.Println(err)
		c.IndentedJSON(http.StatusBadRequest, models.ErrorResponse{Error: "string"})

		return
	}

	var res = database.DB.Db.First(&tagFeatureBanner)
	if res.Error != nil {
		log.Println(res.Error)
		c.IndentedJSON(http.StatusBadRequest, models.ErrorResponse{Error: "string"})

		return
	}

	var banner = models.Banner{}

	banner.ID = tagFeatureBanner.BannerID
	res = database.DB.Db.First(&banner)
	if res.Error != nil {
		log.Println(res.Error)
		c.IndentedJSON(http.StatusInternalServerError, models.ErrorResponse{Error: "string"})

		return
	}

	c.IndentedJSON(http.StatusOK, banner.Content)
}
