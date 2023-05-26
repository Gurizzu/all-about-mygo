package service

import (
	"context"
	"elasticsearch-gin-gonic/src/model"
	"elasticsearch-gin-gonic/src/model/enum"
	db "elasticsearch-gin-gonic/src/util/db/es"
	"fmt"
)

type MovieService struct {
	indexName string
	ctx       context.Context
	esUtil    *db.ElasticUtil
}

func NewMovieService() *MovieService {
	this := &MovieService{
		indexName: enum.ElasticIndex_Movies.String(),
		ctx:       context.Background(),
	}
	this.esUtil = db.NewElasticDbUtilUseEnv(this.indexName)
	return this
}

func (this *MovieService) BaseGetAll(param model.Movies_Search, index *db.ElasticUtil) (data []model.Movies_View,
	metadata model.MetadataResponse) {

	//var filter interface{}
	var listFilterAnd []interface{}
	param.HandleFilter(&listFilterAnd)

	//if len(listFilterAnd) > 0 {
	//	filter = map[string]interface{}{
	//		"query": map[string]interface{}{
	//			"bool": map[string]interface{}{
	//				"must": listFilterAnd,
	//			},
	//		},
	//	}
	//}
	metadata.Pagination, metadata.Message = index.Find(listFilterAnd, param.Request, &data)
	return
}

func (this *MovieService) GetAll(param model.Movies_Search) (data []model.Movies_View, metadata model.MetadataResponse) {
	return this.BaseGetAll(param, this.esUtil)
}

func (this *MovieService) GetOne(key, value string) (res model.Movies_View, errMessage string) {
	query := fmt.Sprintf(`{"size": 1, "query": {"match": {"%s": "%s"}}}`, key, value)
	errMessage = OverrideError(this.esUtil.FindOne(query, &res))
	return
}

func (this *MovieService) Upsert(param model.Movies, isUpdate bool) (resp model.Response) {
	upsertErr, upsertId := this.esUtil.UpsertAndGetId(isUpdate, &param)
	resp.Metadata.Message = upsertErr
	resp.Data = model.Response_Data_Upsert{ID: upsertId}
	return
}

func (this *MovieService) DeleteOne(value string) (errMessage string) {
	fmt.Println(value)
	errMessage = this.esUtil.DeleteOne(value)
	return
}
