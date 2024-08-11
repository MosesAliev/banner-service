package handlers

import (
	"log"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
)

func Auth(next func(c *gin.Context)) func(c *gin.Context) {
	return func(c *gin.Context) {
		defer func() {
			if r := recover(); r != nil {
				log.Println("Recovery:", r)
				c.Status(http.StatusUnauthorized)
			}

		}()

		log.Println(c.GetHeader("Authorization"))
		token := strings.Split(c.GetHeader("Authorization"), " ")[1]

		secretKey := []byte("auth")
		jwtToken, _ := jwt.Parse(token, func(t *jwt.Token) (interface{}, error) {
			return secretKey, nil
		})

		res, ok := jwtToken.Claims.(jwt.MapClaims)
		// обязательно используем второе возвращаемое значение ok и проверяем его, потому что
		// если Сlaims вдруг оказжется другого типа, мы получим панику
		if !ok {
			log.Printf("failed to typecast to jwt.MapCalims")
			return
		}

		roleRaw := res["role"]
		role, ok := roleRaw.(string)
		if !ok {
			log.Printf("failed to typecast to string login")
			return
		}

		log.Println(role)
		c.Request.Header.Add("role", role)
		next(c)
	}

}
