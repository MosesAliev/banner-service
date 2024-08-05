package handlers

import (
	"banner-service/internal/database"
	"banner-service/internal/models"
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// Обновление содержимого баннера
func EditBannerHandler(c *gin.Context) {
	var banner models.Banner
	var err error
	banner.ID, err = strconv.Atoi(c.Param("id"))
	if err != nil {
		log.Println(err)
		c.IndentedJSON(http.StatusBadRequest, models.ErrorResponse{Error: "string"})

		return
	}

	var res *gorm.DB
	res = database.DB.Db.Where("id = ?", banner.ID).First(&models.Banner{}) // Ищем ID баннера в БД

	if res.Error != nil {
		log.Println(err)
		c.IndentedJSON(http.StatusBadRequest, models.ErrorResponse{Error: "string"})

		return
	}

	c.BindJSON(&banner)
	database.DB.Db.Unscoped().Where("banner_id = ?", banner.ID).Delete(&models.TagFeatureBanner{}) // Удаляем устаревшую информацию о баннере

	// Создаем список тегов с фичей баннера
	var tagFeatureBanners []models.TagFeatureBanner
	for _, tagID := range banner.TagIDs {
		banner.Tags = append(banner.Tags, models.Tag{ID: tagID})

		var tagFeatureBanner = models.TagFeatureBanner{}

		tagFeatureBanner.TagID = tagID
		tagFeatureBanner.FeatureID = banner.FeatureID
		tagFeatureBanners = append(tagFeatureBanners, tagFeatureBanner)
	}

	res = database.DB.Db.Save(&banner.Tags) // Добавляем теги баннера в БД
	if res.Error != nil {
		log.Println(res.Error)
		c.IndentedJSON(http.StatusInternalServerError, models.ErrorResponse{Error: "string"})

		return
	}

	res = database.DB.Db.Save(&models.Feature{ID: banner.FeatureID}) // Добавляем фичу баннера в БД
	if res.Error != nil {
		log.Println(res.Error)
		c.IndentedJSON(http.StatusInternalServerError, models.ErrorResponse{Error: "string"})

		return
	}

	var lastAddedFeature = models.Feature{}

	res.Last(&lastAddedFeature)
	banner.FeatureID = lastAddedFeature.ID
	res = database.DB.Db.Save(&banner) // Сохраняем баннер в БД
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

	res = database.DB.Db.Create(&tagFeatureBanners) // Добавляем теги фичи баннеры в БД
	if res.Error != nil {
		log.Println(err)
		c.IndentedJSON(http.StatusInternalServerError, models.ErrorResponse{Error: "string"})

		return
	}

}
