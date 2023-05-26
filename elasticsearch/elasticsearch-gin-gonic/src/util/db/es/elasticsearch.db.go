package es

import (
	"bytes"
	"context"
	"elasticsearch-gin-gonic/src/config"
	"elasticsearch-gin-gonic/src/model"
	"elasticsearch-gin-gonic/src/util"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"
	"reflect"
	"strings"
	"time"

	elastic "github.com/elastic/go-elasticsearch/v8"
	"github.com/elastic/go-elasticsearch/v8/esapi"
)

var Elastic_Live_Connection int

type ElasticUtil struct {
	srv       string
	dbName    string
	indexName string
	ctx       context.Context
}

func NewElasticUtil(srv, indexName string) *ElasticUtil {
	return &ElasticUtil{
		srv:       srv,
		indexName: indexName,
		ctx:       context.Background(),
	}
}

func NewElasticDbUtilUseEnv(indexName string) *ElasticUtil {
	return &ElasticUtil{
		srv:       os.Getenv(config.DB_ELASTIC_SEARCH_SRV),
		indexName: indexName,
		ctx:       context.Background(),
	}
}

func (this ElasticUtil) Connect() (client *elastic.Client, err error) {
	var cfg elastic.Config = elastic.Config{
		Addresses: []string{
			this.srv,
		},
		Username: "",
		Password: "",
	}

	client, err = elastic.NewClient(cfg)
	if err != nil {
		log.Fatal(err)
	}
	return
}

func (this ElasticUtil) UpsertAndGetId(isUpdate bool, ptrParam interface{}) (errMessage string, newDataId string) {
	client, err := this.Connect()
	if err != nil {
		return
	}
	paramAsReflect := reflect.ValueOf(ptrParam).Elem()
	if updatedAt := paramAsReflect.FieldByName("UpdatedAt"); updatedAt.IsValid() {
		updatedAt.SetInt(time.Now().UnixMilli())
	}
	idField := paramAsReflect.FieldByName("IdDocument")

	if !isUpdate {
		if createdAt := paramAsReflect.FieldByName("CreatedAt"); createdAt.IsValid() {
			createdAt.SetInt(time.Now().UnixMilli())
		}
		if idField.IsValid() {
			idField.SetString(util.GenerateID())
		}

		var buf bytes.Buffer
		err := json.NewEncoder(&buf).Encode(ptrParam)
		if err != nil {
			log.Println(err)
		}

		req := esapi.IndexRequest{
			Index:      this.indexName,
			DocumentID: fmt.Sprintf("%s", idField),
			Body:       &buf,
			Refresh:    "true",
		}

		res, err := req.Do(this.ctx, client)
		defer res.Body.Close()
		if err != nil {
			log.Println(err)
			return
		} else {
			var r map[string]interface{}
			if err := json.NewDecoder(res.Body).Decode(&r); err != nil {
				log.Println(err)
			}

			newDataId = r["_id"].(string)
		}
	} else {
		if !idField.IsValid() {
			log.Println("Fail to get field ID")
			return
		}

		bufUpdate, _ := json.Marshal(ptrParam)

		updateReq := esapi.UpdateRequest{
			Index:      this.indexName,
			DocumentID: fmt.Sprintf("%s", idField),
			Body:       bytes.NewReader([]byte(fmt.Sprintf(`{"doc":%s}`, bufUpdate))),
			Refresh:    "true",
		}

		updateRes, err := updateReq.Do(this.ctx, client)

		defer updateRes.Body.Close()
		if err != nil {
			log.Println(err)
			return
		} else {
			var rUpdate map[string]interface{}
			if err := json.NewDecoder(updateRes.Body).Decode(&rUpdate); err != nil {
				log.Println(err)
			}

			newDataId = rUpdate["_id"].(string)

		}

	}
	return
}

func (this ElasticUtil) BaseFindOne(query string, pointerDecodeTo interface{}) (err error) {
	client, err := this.Connect()
	if err != nil {
		return
	}

	res, err := client.Search(
		client.Search.WithContext(this.ctx),
		client.Search.WithIndex(this.indexName),
		client.Search.WithBody(strings.NewReader(query)),
	)

	res, err = client.Search()

	if err != nil {
		log.Println(err)
		return err
	}

	defer res.Body.Close()
	if res.IsError() {
		log.Println(err)
		return err
	}

	var r map[string]interface{}
	if err = json.NewDecoder(res.Body).Decode(&r); err != nil {
		log.Println(err)
		return err
	}

	hits := r["hits"].(map[string]interface{})["hits"].([]interface{})
	if len(hits) == 0 {
		return errors.New("Data not found.")
	}
	hit := hits[0].(map[string]interface{})
	source := hit["_source"]

	jsonString, _ := json.Marshal(source)

	// convert json to struct
	err = json.Unmarshal(jsonString, pointerDecodeTo)
	if err != nil {
		return err
	}
	return
}

func (this ElasticUtil) FindOne(filter string, pointerDecodeTo interface{}) (err error) {
	return this.BaseFindOne(filter, pointerDecodeTo)
}

func (this ElasticUtil) DeleteOne(value string) (errMessage string) {
	client, err := this.Connect()
	if err != nil {
		return
	}

	res, err := client.Delete(
		this.indexName,
		value,
	)
	if err != nil {
		log.Println(err)
		return
	}
	if res.StatusCode == 404 {
		errMessage = fmt.Sprintf("Document with id %s not found", value)
	}

	return
}

func (this ElasticUtil) BaseFindPagination(filter interface{}, fieldRequestParam, pointerDecodeTo interface{}, rangeField string) (paginationResp *model.PaginationResponse, errMessage string) {

	client, err := this.Connect()
	if err != nil {
		return
	}

	var requestPagination model.Request_Pagination
	//* ----------------------------- SET FILTER REQUEST ---------------------------- */
	switch requestAsType := fieldRequestParam.(type) {
	case model.Request:
		filter = requestAsType.BaseHandle(filter, rangeField)
		requestPagination = requestAsType.Request_Pagination
	case model.Request_Pagination:
		requestPagination = requestAsType
	}

	//* ------------------------- SET PAGINATION OPTIONS ------------------------- */
	var ListSort []interface{}
	if orderBy := requestPagination.OrderBy; orderBy != "" {
		sort := map[string]interface{}{
			requestPagination.OrderBy: map[string]interface{}{
				"order": strings.ToLower(requestPagination.Order),
			},
		}
		ListSort = append(ListSort, sort)
	}
	skip, limit := GetSkipAndLimit(requestPagination)

	filter.(map[string]interface{})["skip"] = skip
	filter.(map[string]interface{})["size"] = limit
	filter.(map[string]interface{})["sort"] = ListSort

	jString, _ := json.Marshal(filter)

	res, err := client.Search(
		client.Search.WithContext(this.ctx),
		client.Search.WithIndex(this.indexName),
		client.Search.WithBody(strings.NewReader(fmt.Sprintf(`%s`, jString))),
	)
	defer res.Body.Close()
	if res.IsError() {
		log.Println(err)
		return
	}

	var r map[string]interface{}
	if err := json.NewDecoder(res.Body).Decode(&r); err != nil {
		log.Fatalf("Error parsing the response body: %s", err)
	}

	// totalHits := r["hits"].(map[string]interface{})["total"].(map[string]interface{})["value"]
	paginationResp = &model.PaginationResponse{
		Size: int(limit),
		// TotalElements: int64(totalHits),
		// TotalPages:    int64(math.Ceil(float64(totalHits) / float64(limit))),
	}
	return

}

func (this ElasticUtil) Find(filter interface{}, fieldRequestParam, pointerDecodeTo interface{}) (paginationResponse *model.PaginationResponse, errMessage string) {
	return this.BaseFindPagination(filter, fieldRequestParam, pointerDecodeTo, "")
}
