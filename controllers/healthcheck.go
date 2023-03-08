package controllers

import (
	"net/http"
	"time"

	"byc1/infra"

	"github.com/gin-gonic/gin"
)

func Healthcheck(c *gin.Context) {
	errDbPayment := infra.CheckDB()
	var sDbPayment string
	if errDbPayment != nil {
		sDbPayment = errDbPayment.Error()
	} else {
		sDbPayment = "Coneccion OK"
	}
	if errDbPayment != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"DB_FACTHOSGESTION": sDbPayment,
			"time":              time.Now(),
		})
	} else {
		c.JSON(http.StatusOK, gin.H{
			"DB_FACTHOSGESTION": sDbPayment,
			"time":              time.Now(),
		})
	}
}

func HealthcheckLocal(c *gin.Context) {
	errDbPayment := infra.CheckDBLocal()
	var sDbPayment string
	if errDbPayment != nil {
		sDbPayment = errDbPayment.Error()
	} else {
		sDbPayment = "Coneccion OK"
	}
	if errDbPayment != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"DBPOSTGRES": sDbPayment,
			"time":       time.Now(),
		})
	} else {
		c.JSON(http.StatusOK, gin.H{
			"DBPOSTGRES": sDbPayment,
			"time":       time.Now(),
		})
	}
}
