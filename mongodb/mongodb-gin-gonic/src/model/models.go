package model

type Response struct {
	Data interface{} `json:"data"`
}

type Response_Data_Upsert struct {
	ID string `json:"id"`
}
