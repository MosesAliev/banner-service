package handlers

import (
	"banner-service/internal/database"
	"banner-service/internal/models"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func AddBannerHandler(c *gin.Context) {
	var banner = models.Banner{}

	var err = c.BindJSON(&banner)
	if err != nil {
		log.Println(err)
		c.IndentedJSON(http.StatusBadRequest, models.ErrorResponse{Error: "string"})

		return
	}

	var res *gorm.DB
	var tagFeatureBanner = models.TagFeatureBanner{}

	for _, tagID := range banner.TagIDs {
		banner.Tags = append(banner.Tags, models.Tag{ID: tagID})

		tagFeatureBanner.TagID = tagID
		tagFeatureBanner.FeatureID = banner.FeatureID
		tagFeatureBanner.BannerID = banner.ID
	}

	res = database.DB.Db.Save(&banner.Tags)
	if res.Error != nil {
		log.Println(res.Error)
		c.IndentedJSON(http.StatusInternalServerError, models.ErrorResponse{Error: "string"})

		return
	}

	res = database.DB.Db.Save(&models.Feature{ID: banner.FeatureID})
	if res.Error != nil {
		log.Println(res.Error)
		c.IndentedJSON(http.StatusInternalServerError, models.ErrorResponse{Error: "string"})

		return
	}

	var lastAddedFeature = models.Feature{}

	res.Last(&lastAddedFeature)
	banner.FeatureID = lastAddedFeature.ID
	res = database.DB.Db.Create(&banner)
	if res.Error != nil {
		log.Println(res.Error)
		c.IndentedJSON(http.StatusInternalServerError, models.ErrorResponse{Error: "string"})

		return
	}

	var lastAddedBanner = models.Banner{}

	res.Last(&lastAddedBanner)
	tagFeatureBanner.BannerID = lastAddedBanner.ID
	log.Println(tagFeatureBanner.TagID)
	log.Println(tagFeatureBanner.FeatureID)
	log.Println(tagFeatureBanner.BannerID)
	res = database.DB.Db.Create(&tagFeatureBanner)
	if res.Error != nil {
		log.Println(res.Error)
		c.IndentedJSON(http.StatusBadRequest, models.ErrorResponse{Error: "string"})

		return
	}

	c.IndentedJSON(http.StatusCreated, banner.ID)
}
