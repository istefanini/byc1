package routes

import (
	"byc1/controllers"

	"github.com/gin-gonic/gin"
)

func CreateRoutes(r *gin.Engine) {

	v1 := r.Group("")
	v1.Use()
	{
		v1.GET("/healthcheck", controllers.Healthcheck)
		v1.GET("/healthchecklocal", controllers.HealthcheckLocal)
		v1.GET("/sendfile", controllers.Filehandle)
		v1.POST("/sendfilebyc", controllers.Filehandlebyc)
	}
}
