package main

import (
	"byc1/infra"
	"byc1/routes"
	"os"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
	"github.com/subosito/gotenv"
)

func init() { //connect DB
	_ = gotenv.Load(".env")
}

func main() { //hacerlo con Gin

	infra.SqlConf = &infra.DBData{
		DB_DRIVER:   os.Getenv("DB_DRIVER"),
		DB_USER:     os.Getenv("DB_USERNAME"),
		DB_PASSWORD: os.Getenv("DB_PASSWORD"),
		DB_HOST:     os.Getenv("DB_HOST"),
		DB_INSTANCE: os.Getenv("DB_INSTANCE"),
		DB_DATABASE: os.Getenv("DB_DATABASE"),
		DB_ENCRYPT:  os.Getenv("DB_ENCRYPT"),
	}
	infra.DbPayment = infra.ConnectDB()
	defer infra.DbPayment.Close()

	infra.DbLocal = infra.ConnectDBLocal()
	defer infra.DbLocal.Close()

	gin.SetMode(gin.ReleaseMode)
	r := gin.Default()
	r.Use(cors.Default())
	routes.CreateRoutes(r)
	serverPort := os.Getenv("API_PORT")
	_ = r.Run(":" + serverPort)
}
