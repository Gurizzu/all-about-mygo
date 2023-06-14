package main

import (
	"fmt"
	"log"
	"mongodb-gin-gonic/docs"
	"mongodb-gin-gonic/src/controller"

	"github.com/gin-gonic/gin"
	"github.com/logrusorgru/aurora"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func main() {
	appPort := ":13456"

	router := gin.Default()

	basePath := "/api/v1"
	docs.SwaggerInfo.BasePath = basePath
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	apiV1 := router.Group(basePath)

	log.Println(aurora.Green(
		fmt.Sprintf("http://localhost%s/swagger/index.html", appPort),
	))
	controller.NewItemController(apiV1)
	log.Fatalln(router.Run(appPort))

}
