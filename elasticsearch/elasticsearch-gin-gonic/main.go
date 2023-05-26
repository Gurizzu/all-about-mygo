package main

import (
	"elasticsearch-gin-gonic/docs"
	"elasticsearch-gin-gonic/src/config"
	"elasticsearch-gin-gonic/src/controller"
	"fmt"

	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/gzip"
	"github.com/gin-gonic/gin"
	"github.com/logrusorgru/aurora"
	"github.com/subosito/gotenv"

	"log"
	"os"
	"time"

	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func init() {
	log.SetFlags(log.LstdFlags | log.Llongfile)
	if err := gotenv.Load(".env"); err != nil {
		log.Printf("env Load Error : %v", err)
	}
	fmt.Println("Using timezone:", aurora.Green(time.Now().Location().String()))
}

func main() {
	appPort := ":" + os.Getenv(config.ENV_KEY_PORT)
	fmt.Printf("ini port = %s", appPort)

	if os.Getenv("LOCAL") != "true" {
		gin.SetMode(gin.ReleaseMode)
	}
	router := gin.Default()

	router.Use(gzip.Gzip(gzip.BestSpeed))
	router.Use(cors.Default())

	basePath := "/api/v1"
	docs.SwaggerInfo.BasePath = basePath
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	apiV1 := router.Group(basePath)
	apiV1.GET("/ping", func(c *gin.Context) { c.JSON(200, gin.H{"message": "Nuwhofev"}) })

	controller.NewMoviesController(apiV1)

	log.Println(aurora.Green(
		fmt.Sprintf("http://localhost%s/swagger/index.html", appPort),
	))
	log.Fatalln(router.Run(appPort))
}
