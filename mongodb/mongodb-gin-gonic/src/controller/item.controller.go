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
}
