package model

type Item struct {
	Id          string `json:"id" bson:"_id"`
	Name        string `json:"name" bson:"name"`
	Description string `json:"description" bson:"description"`
	Created_At  int64  `json:"created_at" bson:"created_at"`
	Updated_At  int64  `json:"updated_at" bson:"updated_at"`
}
