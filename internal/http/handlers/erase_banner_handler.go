package handlers

import (
	"banner-service/internal/database"
	"banner-service/internal/models"
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// Удаление баннера по идентификатору
func EraseBannerHandler(c *gin.Context) {
	var ID, err = strconv.Atoi(c.Param("id"))
	if err != nil {
		log.Println(err)
		c.IndentedJSON(http.StatusBadRequest, models.ErrorResponse{Error: "string"})

		return
	}

	var res = database.DB.Db.Where("id = ?", ID).First(&models.Banner{})

	if res.Error != nil {
		log.Println(err)
		c.IndentedJSON(http.StatusBadRequest, models.ErrorResponse{Error: "string"})

		return
	}

	res = database.DB.Db.Unscoped().Where("banner_id = ?", ID).Delete(&models.TagFeatureBanner{})

	if res.Error != nil {
		log.Println(err)
		c.IndentedJSON(http.StatusBadRequest, models.ErrorResponse{Error: "string"})

		return

	}

	res = database.DB.Db.Unscoped().Where("id = ?", ID).Delete(&models.Banner{})
	if res.Error != nil {
		log.Println(err)
		c.IndentedJSON(http.StatusInternalServerError, models.ErrorResponse{Error: "string"})

		return
	}

	c.Status(http.StatusNoContent)
}
