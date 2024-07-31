package handlers

import (
	"banner-service/internal/database"
	"banner-service/internal/models"
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func UserBannersHanlder(c *gin.Context) {
	var banners []models.Banner
	var strTagID, tagOK = c.GetQuery("tag_id")
	var strFeatureID, featureOK = c.GetQuery("feature_id")
	var strLimit, limitOK = c.GetQuery("limit")
	var strOffset, offsetOK = c.GetQuery("offset")
	if tagOK && featureOK && limitOK && offsetOK {
		var tagFeatureBanners []models.TagFeatureBanner

		var err error
		var tagID int
		tagID, err = strconv.Atoi(strTagID)
		if err != nil {
			log.Println(err)
			c.IndentedJSON(http.StatusBadRequest, models.ErrorResponse{Error: "string"})
			return
		}

		var featureID int
		featureID, err = strconv.Atoi(strFeatureID)
		if err != nil {
			log.Println(err)
			c.IndentedJSON(http.StatusBadRequest, models.ErrorResponse{Error: "string"})
			return
		}

		var limit int
		limit, err = strconv.Atoi(strLimit)
		if err != nil {
			log.Println(err)
			c.IndentedJSON(http.StatusBadRequest, models.ErrorResponse{Error: "string"})
			return
		}

		var offset int
		offset, err = strconv.Atoi(strOffset)
		if err != nil {
			log.Println(err)
			c.IndentedJSON(http.StatusBadRequest, models.ErrorResponse{Error: "string"})
			return
		}

		var res = database.DB.Db.Limit(limit).Offset(offset).Where("tag_id = ?", tagID).Or("feature_id = ?", featureID).Find(&tagFeatureBanners)
		if res.Error != nil {
			log.Println(res.Error)
			c.IndentedJSON(http.StatusInternalServerError, models.ErrorResponse{Error: "string"})
			return
		}

		for _, tagFeatureBanner := range tagFeatureBanners {
			var banner = models.Banner{}

			res = database.DB.Db.Where("id = ?", tagFeatureBanner.BannerID).First(&banner)
			if res.Error != nil {
				log.Println(res.Error)
				c.IndentedJSON(http.StatusInternalServerError, models.ErrorResponse{Error: "string"})
				return
			}

			var thisBannerTagFeature []models.TagFeatureBanner
			res = database.DB.Db.Where("banner_id = ?", banner.ID).Find(&thisBannerTagFeature)
			if res.Error != nil {
				log.Println(res.Error)
				c.IndentedJSON(http.StatusInternalServerError, models.ErrorResponse{Error: "string"})
				return
			}

			for i := 0; i < len(tagFeatureBanners); i++ {
				banner.TagIDs = append(banner.TagIDs, tagFeatureBanner.TagID)
			}

			banners = append(banners, banner)
		}

		c.IndentedJSON(http.StatusOK, banners)
		return
	}

	var res = database.DB.Db.Find(&banners)
	if res.Error != nil {
		log.Println(res.Error)
		c.IndentedJSON(http.StatusInternalServerError, models.ErrorResponse{Error: "string"})
		return
	}

	c.IndentedJSON(http.StatusOK, banners)
}
