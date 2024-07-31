package main

import (
	"banner-service/internal/database"
	"banner-service/internal/http/router"
)

func main() {
	database.ConnectDB()
	var r = router.SetupRouter()
	r.Run()
}
