package model

type Response struct {
	Metadata MetadataResponse `bson:",inline"`
	Data     interface{}      `json:"data"`
}

type MetadataResponse struct {
	Status  bool   `json:"status"`
	Message string `json:"message"`
}

type Response_Data_Upsert struct {
	ID string `json:"id"`
}
