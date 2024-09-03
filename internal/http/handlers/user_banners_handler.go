package handlers

import (
	"banner-service/internal/database"
	"banner-service/internal/models"
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// Получение всех баннеров c фильтрацией по фиче и/или тегу
func UserBannersHanlder(c *gin.Context) {
	// проверка прав доступа
	var role = c.GetHeader("role")
	if role != "user" && role != "admin" {
		c.Status(http.StatusForbidden)
		return
	}

	var banners []models.Banner
	var strTagID, tagOK = c.GetQuery("tag_id")             // Получаем параметры для фильтрации
	var strFeatureID, featureOK = c.GetQuery("feature_id") //
	var strLimit, limitOK = c.GetQuery("limit")            //
	var strOffset, offsetOK = c.GetQuery("offset")         //
	if tagOK && featureOK && limitOK && offsetOK {
		var tagFeatureBanners []models.TagFeatureBanner // Список тегов фич и баннеров
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

		// Применяем параметры для поиска
		var res = database.DB.Db.Limit(limit).Offset(offset).Where("tag_id = ?", tagID).Or("feature_id = ?", featureID).Find(&tagFeatureBanners)
		if res.Error != nil {
			log.Println(res.Error)
			c.IndentedJSON(http.StatusInternalServerError, models.ErrorResponse{Error: "string"})

			return
		}

		// Перебираем полученный список
		for _, tagFeatureBanner := range tagFeatureBanners {
			var banner = models.Banner{}

			res = database.DB.Db.Where("id = ?", tagFeatureBanner.BannerID).First(&banner)
			if res.Error != nil {
				log.Println(res.Error)
				c.IndentedJSON(http.StatusInternalServerError, models.ErrorResponse{Error: "string"})

				return
			}

			var thisBannerTagFeature []models.TagFeatureBanner // список всех тегов и фич конкретного баннера
			res = database.DB.Db.Where("banner_id = ?", banner.ID).Find(&thisBannerTagFeature)
			if res.Error != nil {
				log.Println(res.Error)
				c.IndentedJSON(http.StatusInternalServerError, models.ErrorResponse{Error: "string"})

				return
			}

			for i := 0; i < len(thisBannerTagFeature); i++ {
				banner.TagIDs = append(banner.TagIDs, thisBannerTagFeature[i].TagID) // добавляем в список тегов баннера идентификаторы тегов
			}

			banners = append(banners, banner)
		}

		c.IndentedJSON(http.StatusOK, banners) // записываем в тело ответа список полученных баннеров
		return
	}

	if tagOK && featureOK && offsetOK {
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

		var offset int
		offset, err = strconv.Atoi(strOffset)
		if err != nil {
			log.Println(err)
			c.IndentedJSON(http.StatusBadRequest, models.ErrorResponse{Error: "string"})

			return
		}

		var res = database.DB.Db.Offset(offset).Where("tag_id = ?", tagID).Or("feature_id = ?", featureID).Find(&tagFeatureBanners)
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

			for i := 0; i < len(thisBannerTagFeature); i++ {
				banner.TagIDs = append(banner.TagIDs, thisBannerTagFeature[i].TagID)
			}

			banners = append(banners, banner)
		}

		c.IndentedJSON(http.StatusOK, banners)
		return
	}

	if tagOK && featureOK && limitOK {
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

		var res = database.DB.Db.Limit(limit).Where("tag_id = ?", tagID).Or("feature_id = ?", featureID).Find(&tagFeatureBanners)
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

			for i := 0; i < len(thisBannerTagFeature); i++ {
				banner.TagIDs = append(banner.TagIDs, thisBannerTagFeature[i].TagID)
			}

			banners = append(banners, banner)
		}

		c.IndentedJSON(http.StatusOK, banners)
		return
	}

	if tagOK && featureOK {
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

		var res = database.DB.Db.Where("tag_id = ?", tagID).Or("feature_id = ?", featureID).Find(&tagFeatureBanners)
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

			for i := 0; i < len(thisBannerTagFeature); i++ {
				banner.TagIDs = append(banner.TagIDs, thisBannerTagFeature[i].TagID)
			}

			banners = append(banners, banner)
		}

		c.IndentedJSON(http.StatusOK, banners)
		return
	}

	if tagOK && limitOK {
		var tagFeatureBanners []models.TagFeatureBanner
		var err error
		var tagID int
		tagID, err = strconv.Atoi(strTagID)
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

		var res = database.DB.Db.Limit(limit).Where("tag_id = ?", tagID).Find(&tagFeatureBanners)
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

			for i := 0; i < len(thisBannerTagFeature); i++ {
				banner.TagIDs = append(banner.TagIDs, thisBannerTagFeature[i].TagID)
			}

			banners = append(banners, banner)
		}

		c.IndentedJSON(http.StatusOK, banners)
		return
	}

	if tagOK && offsetOK {
		var tagFeatureBanners []models.TagFeatureBanner
		var err error
		var tagID int
		tagID, err = strconv.Atoi(strTagID)
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

		var res = database.DB.Db.Offset(offset).Where("tag_id = ?", tagID).Find(&tagFeatureBanners)
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

			for i := 0; i < len(thisBannerTagFeature); i++ {
				banner.TagIDs = append(banner.TagIDs, thisBannerTagFeature[i].TagID)
			}

			banners = append(banners, banner)
		}

		c.IndentedJSON(http.StatusOK, banners)
		return
	}

	if featureOK && limitOK {
		var tagFeatureBanners []models.TagFeatureBanner
		var err error
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

		var res = database.DB.Db.Limit(limit).Where("feature_id = ?", featureID).Find(&tagFeatureBanners)
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

			for i := 0; i < len(thisBannerTagFeature); i++ {
				banner.TagIDs = append(banner.TagIDs, thisBannerTagFeature[i].TagID)
			}

			banners = append(banners, banner)
		}

		c.IndentedJSON(http.StatusOK, banners)
		return
	}

	if featureOK && offsetOK {
		var tagFeatureBanners []models.TagFeatureBanner
		var err error
		var featureID int
		featureID, err = strconv.Atoi(strFeatureID)
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

		var res = database.DB.Db.Offset(offset).Where("feature_id = ?", featureID).Find(&tagFeatureBanners)
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

			for i := 0; i < len(thisBannerTagFeature); i++ {
				banner.TagIDs = append(banner.TagIDs, thisBannerTagFeature[i].TagID)
			}

			banners = append(banners, banner)
		}

		c.IndentedJSON(http.StatusOK, banners)
		return
	}

	if tagOK {
		var tagFeatureBanners []models.TagFeatureBanner
		var err error
		var tagID int
		tagID, err = strconv.Atoi(strTagID)
		if err != nil {
			log.Println(err)
			c.IndentedJSON(http.StatusBadRequest, models.ErrorResponse{Error: "string"})

			return
		}

		var res = database.DB.Db.Where("tag_id = ?", tagID).Find(&tagFeatureBanners)
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

			for i := 0; i < len(thisBannerTagFeature); i++ {
				banner.TagIDs = append(banner.TagIDs, thisBannerTagFeature[i].TagID)
			}

			banners = append(banners, banner)
		}

		c.IndentedJSON(http.StatusOK, banners)
		return
	}

	if featureOK {
		var tagFeatureBanners []models.TagFeatureBanner
		var err error
		var featureID int
		featureID, err = strconv.Atoi(strFeatureID)
		if err != nil {
			log.Println(err)
			c.IndentedJSON(http.StatusBadRequest, models.ErrorResponse{Error: "string"})

			return
		}

		var res = database.DB.Db.Where("feature_id = ?", featureID).Find(&tagFeatureBanners)
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

			for i := 0; i < len(thisBannerTagFeature); i++ {
				banner.TagIDs = append(banner.TagIDs, thisBannerTagFeature[i].TagID)
			}

			banners = append(banners, banner)
		}

		c.IndentedJSON(http.StatusOK, banners)
		return
	}

	if limitOK && offsetOK {
		var tagFeatureBanners []models.TagFeatureBanner
		var err error

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

		var res = database.DB.Db.Limit(limit).Offset(offset).Find(&tagFeatureBanners)
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

			for i := 0; i < len(thisBannerTagFeature); i++ {
				banner.TagIDs = append(banner.TagIDs, thisBannerTagFeature[i].TagID)
			}

			banners = append(banners, banner)
		}

		c.IndentedJSON(http.StatusOK, banners)
		return
	}

	if limitOK {
		var tagFeatureBanners []models.TagFeatureBanner
		var err error
		var limit int
		limit, err = strconv.Atoi(strLimit)
		if err != nil {
			log.Println(err)
			c.IndentedJSON(http.StatusBadRequest, models.ErrorResponse{Error: "string"})

			return
		}

		var res = database.DB.Db.Limit(limit).Find(&tagFeatureBanners)
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

			for i := 0; i < len(thisBannerTagFeature); i++ {
				banner.TagIDs = append(banner.TagIDs, thisBannerTagFeature[i].TagID)
			}

			banners = append(banners, banner)
		}

		c.IndentedJSON(http.StatusOK, banners)
		return
	}

	if offsetOK {
		var tagFeatureBanners []models.TagFeatureBanner
		var err error
		var offset int
		offset, err = strconv.Atoi(strOffset)
		if err != nil {
			log.Println(err)
			c.IndentedJSON(http.StatusBadRequest, models.ErrorResponse{Error: "string"})

			return
		}

		var res = database.DB.Db.Offset(offset).Find(&tagFeatureBanners)
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

			for i := 0; i < len(thisBannerTagFeature); i++ {
				banner.TagIDs = append(banner.TagIDs, thisBannerTagFeature[i].TagID)
			}

			banners = append(banners, banner)
		}

		c.IndentedJSON(http.StatusOK, banners)
		return
	}

	var res = database.DB.Db.Find(&banners)
	for i := 0; i < len(banners); i++ {
		var thisBannerTagFeature []models.TagFeatureBanner
		res = database.DB.Db.Where("banner_id = ?", banners[i].ID).Find(&thisBannerTagFeature)
		if res.Error != nil {
			log.Println(res.Error)
			c.IndentedJSON(http.StatusInternalServerError, models.ErrorResponse{Error: "string"})

			return
		}

		for j := 0; j < len(thisBannerTagFeature); j++ {
			banners[i].TagIDs = append(banners[i].TagIDs, thisBannerTagFeature[j].TagID)
		}

	}

	if res.Error != nil {
		log.Println(res.Error)
		c.IndentedJSON(http.StatusInternalServerError, models.ErrorResponse{Error: "string"})

		return
	}

	c.IndentedJSON(http.StatusOK, banners)
}
