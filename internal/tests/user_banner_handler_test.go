package tests

import (
	"banner-service/internal/database"
	"banner-service/internal/http/router"
	"banner-service/internal/models"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/go-playground/assert"
	"github.com/golang-jwt/jwt"
)

func TestUserBannerHandler(t *testing.T) {
	database.ConnectDB()
	database.DB.Db.Exec("DELETE FROM tag_feature_banners") // Очистка БД
	database.DB.Db.Exec("DELETE FROM tag_ids")             //
	database.DB.Db.Exec("DELETE FROM tags")                //
	database.DB.Db.Exec("DELETE FROM banners")             //
	database.DB.Db.Exec("DELETE FROM features")            //
	// Заполняем БД
	var content = models.Content{}
	content.Title = "some_title"
	content.Text = "some_text"
	content.Url = "some_url"
	var banner = models.Banner{}

	banner.ID = 1
	var tag = models.Tag{}

	tag.ID = 1
	database.DB.Db.Create(&tag)
	banner.Tags = append(banner.Tags, tag)
	banner.FeatureID = 1
	var feature = models.Feature{}

	feature.ID = banner.FeatureID
	database.DB.Db.Create(&feature)
	banner.Content = content
	banner.IsActive = true
	database.DB.Db.Create(&banner)
	var tagFeatureBanner = models.TagFeatureBanner{}
	tagFeatureBanner.BannerID = 1
	tagFeatureBanner.TagID = 1
	tagFeatureBanner.FeatureID = 1
	database.DB.Db.Create(&tagFeatureBanner)
	var r = router.SetupRouter()                                                 // Проверка авторизации
	var w = httptest.NewRecorder()                                               //
	var req, _ = http.NewRequest("GET", "/user_banner", nil)                     //
	req.URL.RawQuery = url.Values{"tag_id": {"1"}, "feature_id": {"1"}}.Encode() //
	//
	r.ServeHTTP(w, req)          //
	assert.Equal(t, 401, w.Code) //
	payload := jwt.MapClaims{
		"role": "user",
	}

	// Позитивный сценарий
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, payload)
	secretKey := []byte("auth")
	signedToken, _ := token.SignedString(secretKey)
	w = httptest.NewRecorder()
	req.Header.Add("Authorization", fmt.Sprintf("token %s", signedToken))
	r.ServeHTTP(w, req)
	assert.Equal(t, 200, w.Code)
	// Не найдены баннеры
	w = httptest.NewRecorder()
	req.URL.RawQuery = url.Values{"tag_id": {"2"}, "feature_id": {"2"}}.Encode()

	r.ServeHTTP(w, req)
	assert.Equal(t, 404, w.Code)
	// Неправильно составлен запрос
	req.URL.RawQuery = url.Values{"tag_id": {"a"}, "feature_id": {"a"}}.Encode()

	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)
	assert.Equal(t, 400, w.Code)
	payload = jwt.MapClaims{
		"role": "use",
	}

	token = jwt.NewWithClaims(jwt.SigningMethodHS256, payload)
	secretKey = []byte("auth")
	signedToken, _ = token.SignedString(secretKey)
	req, _ = http.NewRequest("GET", "/user_banner", nil)
	req.Header.Add("Authorization", fmt.Sprintf("token %s", signedToken))
	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)
	assert.Equal(t, 403, w.Code)

}
