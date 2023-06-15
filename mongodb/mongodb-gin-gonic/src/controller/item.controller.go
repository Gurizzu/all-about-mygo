package controller

import (
	"log"
	"mongodb-gin-gonic/src/model"
	"mongodb-gin-gonic/src/services"

	"github.com/gin-gonic/gin"
)

type ItemController struct {
	router  *gin.RouterGroup
	service *services.ItemService
}

func NewItemController(router *gin.RouterGroup) *ItemController {
	o := &ItemController{router: router, service: services.NewItemService()}

	item := o.router.Group("/item")
	item.POST("/add", o.Add)
	return o
}

// @Tags Item
// @Accept json
// @Param parameter body model.Item true "PARAM"
// @Produce json
// @Success 201 {object} object{meta_data=model.Response} "OK"
// @Router /item/add [post]
// @Security JWT
func (o *ItemController) Add(ctx *gin.Context) {
	resp := model.Response{}

	var param model.Item
	if err := ctx.BindJSON(&param); err != nil {
		log.Println(err)
		return
	}

	resp.Data = o.service.Upsert(param, false)
	ctx.JSON(200, resp)
}

// @Tags Item
// @Accept json
// @Param parameter body model.Item true "PARAM"
// @Produce json
// @Success 200 {object} object{meta_data=model.Response} "OK"
// @Router /item/update [put]
// @Security JWT
func (o *ItemController) Update(ctx *gin.Context) {
	resp := model.Response{}

	var param model.Item
	if err := ctx.BindJSON(&param); err != nil {
		log.Println(err)
		return
	}

	resp.Data = o.service.Upsert(param, true)
	ctx.JSON(200, resp)
}

// @Tags OPD
// @Accept json
// @Param id query string true "ID"
// @Produce json
// @Success 200 {object} object{data=model.OPD_View,meta_data=model.MetadataResponse} "OK"
// @Router /opd/get-one [get]
// @Security JWT
func (o *ItemController) FindOne(ctx *gin.Context) {
	resp := model.Response{}
	resp.Data, resp.Metadata.Message = o.service.FindOne("_id", ctx.Query("id"))
	ctx.JSON(200, resp)
}
