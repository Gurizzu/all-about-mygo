package controller

import (
	"elasticsearch-gin-gonic/src/model"
	"elasticsearch-gin-gonic/src/service"
	"log"
	"time"

	"github.com/gin-gonic/gin"
)

type MoviesController struct {
	router  *gin.RouterGroup
	service *service.MovieService
}

func NewMoviesController(router *gin.RouterGroup) *MoviesController {
	this := &MoviesController{router: router, service: service.NewMovieService()}

	movie := this.router.Group("/movie")
	movie.POST("/get-all", this.GetAll)
	movie.GET("/get-one", this.GetOne)
	movie.POST("/add", this.Add)
	movie.PUT("/update", this.Update)
	movie.DELETE("/delete", this.DeleteOne)

	return this
}

// @Tags movies
// @Accept json
// @Param parameter body model.Movies_Search true "PARAM"
// @Produce json
// @Success 200 {object} object{data=[]model.Movies_View,meta_data=model.MetadataResponse} "OK"
// @Router /movie/get-all [post]
func (this *MoviesController) GetAll(ctx *gin.Context) {
	resp := model.Response{}
	defer SetMetadataResponse(ctx, time.Now(), &resp)

	var param model.Movies_Search
	if err := ctx.BindJSON(&param); err != nil {
		log.Println(err)
		return
	}

	resp.Data, resp.Metadata = this.service.GetAll(param)
}

// @Tags movies
// @Accept json
// @Param id query string true "ID"
// @Produce json
// @Success 200 {object} object{data=model.Movies_View,meta_data=model.MetadataResponse} "OK"
// @Router /movie/get-one [get]
func (this *MoviesController) GetOne(ctx *gin.Context) {
	resp := model.Response{}
	defer SetMetadataResponse(ctx, time.Now(), &resp)

	resp.Data, resp.Metadata.Message = this.service.GetOne("_id", ctx.Query("id"))
}

// @Tags movies
// @Accept json
// @Param parameter body model.Movies true "PARAM"
// @Produce json
// @Success 201 {object} object{meta_data=model.MetadataResponse} "OK"
// @Router /movie/add [post]
func (this *MoviesController) Add(ctx *gin.Context) {
	resp := model.Response{}
	defer SetMetadataResponse(ctx, time.Now(), &resp)

	var param model.Movies
	if err := ctx.BindJSON(&param); err != nil {
		log.Println(err)
		return
	}

	resp = this.service.Upsert(param, false)
}

// @Tags movies
// @Accept json
// @Param parameter body model.Movies true "PARAM"
// @Produce json
// @Success 200 {object} object{meta_data=model.MetadataResponse} "OK"
// @Router /movie/update [put]
func (this *MoviesController) Update(ctx *gin.Context) {
	resp := model.Response{}
	defer SetMetadataResponse(ctx, time.Now(), &resp)

	var param model.Movies
	if err := ctx.BindJSON(&param); err != nil {
		log.Println(err)
		return
	}

	resp = this.service.Upsert(param, true)
}

// @Tags movies
// @Accept json
// @Param id query string true "ID"
// @Produce json
// @Success 200 {object} object{meta_data=model.MetadataResponse} "OK"
// @Router /movie/delete [delete]
// @Security JWT
func (this *MoviesController) DeleteOne(ctx *gin.Context) {
	resp := model.Response{}
	defer SetMetadataResponse(ctx, time.Now(), &resp)
	resp.Metadata.Message = this.service.DeleteOne(ctx.Query("id"))
}
