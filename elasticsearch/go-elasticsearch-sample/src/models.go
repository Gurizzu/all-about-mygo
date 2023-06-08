package src

type Item struct {
	Id          string `json:"id"`
	Title       string `json:"name"`
	Description string `json:"description"`
	Stock       int    `json:"stock"`
	Status      string `json:"status"`
	Category    string `json:"category"`
	CreatedAt   int64  `json:"createdAt,omitempty" bson:"createdAt,omitempty"`
	UpdatedAt   int64  `json:"updatedAt" bson:"updatedAt"`
}

type ItemSearch struct {
	Title  string `json:"title"`
	Type   string `json:"type"`
	Status string `json:"status"`
}

func (o *ItemSearch) HandleFilter(listFilterAnd *[]interface{}) {}
