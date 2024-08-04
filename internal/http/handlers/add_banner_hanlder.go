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
	var tagFeatureBanners []models.TagFeatureBanner
	for _, tagID := range banner.TagIDs {
		banner.Tags = append(banner.Tags, models.Tag{ID: tagID})

		var tagFeatureBanner = models.TagFeatureBanner{}

		tagFeatureBanner.TagID = tagID
		tagFeatureBanner.FeatureID = banner.FeatureID
		tagFeatureBanners = append(tagFeatureBanners, tagFeatureBanner)
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
	for i := 0; i < len(tagFeatureBanners); i++ {
		tagFeatureBanners[i].BannerID = lastAddedBanner.ID
	}

	res = database.DB.Db.Create(&tagFeatureBanners)
	if res.Error != nil {
		log.Println(res.Error)
		c.IndentedJSON(http.StatusBadRequest, models.ErrorResponse{Error: "string"})

		return
	}

	c.IndentedJSON(http.StatusCreated, banner.ID)
}
