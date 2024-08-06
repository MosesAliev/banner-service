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
	var role = c.GetHeader("role")
	if role != "user" {
		c.Status(http.StatusForbidden)
		return
	}

	var tagID, tagOk = c.GetQuery("tag_id") // проверяем наличие тега в запросе и получаем его
	if !tagOk {
		log.Println("Отсутствует тег в запросе")
		c.IndentedJSON(http.StatusBadRequest, models.ErrorResponse{Error: "string"})

		return
	}

	var featureID, featureOk = c.GetQuery("feature_id") // проверяем наличие фичи в запросе и получаем его
	if !featureOk {
		log.Println("Отсутствует фича в запросе")
		c.IndentedJSON(http.StatusBadRequest, models.ErrorResponse{Error: "string"})

		return
	}

	var tagFeatureBanner = models.TagFeatureBanner{} // объект содержит запись тега и фичи для однозначного определения баннера

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

	var res = database.DB.Db.First(&tagFeatureBanner) // находим ID баннера
	if res.Error != nil {
		log.Println(res.Error)
		c.IndentedJSON(http.StatusBadRequest, models.ErrorResponse{Error: "string"})

		return
	}

	var banner = models.Banner{}

	banner.ID = tagFeatureBanner.BannerID
	res = database.DB.Db.First(&banner) // Ищем информацию о баннере в БД
	if res.Error != nil {
		log.Println(res.Error)
		c.IndentedJSON(http.StatusInternalServerError, models.ErrorResponse{Error: "string"})

		return
	}

	c.IndentedJSON(http.StatusOK, banner.Content) // Пишем в тело ответа содержимое баннера
}
