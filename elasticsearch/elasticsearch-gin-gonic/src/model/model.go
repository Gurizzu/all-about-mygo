package model

import "time"

type Metadata struct {
	CreatedAt int64 `json:"createdAt,omitempty" bson:"createdAt,omitempty"`
	UpdatedAt int64 `json:"updatedAt" bson:"updatedAt"`
}

type MetadataWithID struct {
	IdDocument string `json:"id"`
	Metadata
}

type Request_Search struct {
	Range *Range `json:"range"`
}

func (this Request_Search) BaseHandle(filter interface{}, rangeField string) (res interface{}) {
	if requestRange := this.Range; requestRange != nil && requestRange.Field != "" {
		rangeField = this.Range.Field
	}
	if rangeField == "" {
		rangeField = "updatedAt"
	}
	res = filter

	if this.Range != nil {
		if this.Range.Start == 0 && this.Range.End == 0 {
			timeNow := time.Now()
			this.Range.End = timeNow.UnixMilli()
			this.Range.Start = timeNow.AddDate(0, 0, -7).UnixMilli()
		}
		filter = map[string]interface{}{
			rangeField: map[string]interface{}{
				"gte": this.Range.Start,
				"lte": this.Range.End,
			},
		}
	}

	return

}

func (this Request_Search) Handle_RequestSearch(filter interface{}) (res interface{}) {
	return this.BaseHandle(filter, "")
}

type Range struct {
	Field string `json:"field" example:"updatedAt"`
	Start int64  `json:"start" example:"1646792565000"`
	End   int64  `json:"end" example:"1646792565000"`
}

type MetadataResponse struct {
	Status        bool   `json:"status"`
	Message       string `json:"message"`
	TimeExecution string `json:"timeExecution"`

	Pagination *PaginationResponse `json:"pagination" bson:"-"`
}

type Response struct {
	Metadata MetadataResponse `json:"metaData"`
	Data     interface{}      `json:"data"`
}

type PaginationResponse struct {
	Size          int   `json:"size"`
	TotalPages    int64 `json:"totalPages"`
	TotalElements int64 `json:"totalElements"`
}

type Request struct {
	Request_Pagination
	Request_Search
}

type Request_Pagination struct {
	OrderBy string `json:"orderBy" example:"createdAt"`
	Order   string `json:"order" example:"ASC"`

	Page int64 `example:"1" json:"page"`
	Size int64 `example:"11" json:"size"`
}

type Response_Data_Upsert struct {
	ID string `json:"id"`
}
