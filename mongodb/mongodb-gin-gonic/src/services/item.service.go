package services

import (
	"context"
	"errors"
	"log"
	"mongodb-gin-gonic/src/model"
	"mongodb-gin-gonic/src/util/db"
	"mongodb-gin-gonic/src/util/db/ustring"
	"mongodb-gin-gonic/src/util/enum"

	"go.mongodb.org/mongo-driver/bson"
)

type ItemService struct {
	collectionName string
	ctx            context.Context
	dbUtil         *db.MongoDbUtil
}

func NewItemService() *ItemService {
	o := &ItemService{
		collectionName: enum.MongoCollection_Items.String(),
		ctx:            context.Background(),
	}

	o.dbUtil = db.NewMongoDbUtil("mongodb://localhost:27017", "redis-practice", enum.MongoCollection_Items.String())
	return o
}

func (o *ItemService) Upsert(param model.Item, isUpdate bool) (resp model.Response) {

	_, col := o.dbUtil.GetCollection()

	paramId := param.Id
	if paramId == "" && len(paramId) == 0 {
		param.Id = ustring.GenerateID()
	}

	if !isUpdate {
		res, err := col.InsertOne(o.ctx, param)
		if err != nil {
			log.Println(err)
		}
		resp.Data = res.InsertedID
		return
	} else {
		if updateRes, err := col.UpdateByID(o.ctx, bson.M{"_id": param.Id}, bson.M{"$set": param}); err != nil {
			log.Println(err)
		} else {
			if updateRes.MatchedCount == 0 && updateRes.UpsertedID == "" {
				err = errors.New("data not found. nothing updated")
				log.Println(err)
				return
			}
			resp.Data = updateRes.UpsertedID
		}

	}
	return
}

func (o *ItemService) FindOne(key string, value string) (pointerDecodeTo model.Item, err string) {
	_, col := o.dbUtil.GetCollection()

	filter := bson.M{key: value}
	res := col.FindOne(o.ctx, filter)
	if res.Err() != nil {
		log.Println(res.Err(), filter)
		err = "data not found"
		return
	}

	if err := res.Decode(&pointerDecodeTo); err != nil {
		log.Println(err)
	}
	return
}
