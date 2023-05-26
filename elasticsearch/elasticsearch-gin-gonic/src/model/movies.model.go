package model

type Movies struct {
	MetadataWithID `bson:",inline"`
	Duration       string `json:"duration"`
	ListedIn       string `json:"listed_in"`
	Country        string `json:"country"`
	DateAdded      string `json:"date_added"`
	ShowId         string `json:"show_id"`
	Director       string `json:"director"`
	ReleaseYear    int    `json:"release_year"`
	Rating         string `json:"rating"`
	Description    string `json:"description"`
	Types          string `json:"types"`
	Title          string `json:"title"`
}

type Movies_View struct {
	Movies `bson:",inline"`
}

type Movies_Search struct {
	Title string   `json:"search"`
	Genre []string `json:"genre"`

	Request
}

func (this *Movies_Search) HandleFilter(listFilterAnd *[]interface{}) {
	if search := this.Title; search != "" {
		//a := gin.H{ "match_phrase_prefix" : gin.H{"title" : search}}
		var matchQuery = map[string]interface{}{
			"match_phrase_prefix": map[string]interface{}{
				"title": search,
			},
		}
		*listFilterAnd = append(*listFilterAnd, matchQuery)
	}

	if search := this.Genre; len(search) != 0 {
		for _, searchFilter := range search {
			var matchQuery = map[string]interface{}{
				"match": map[string]interface{}{
					"listed_in": searchFilter,
				},
			}
			*listFilterAnd = append(*listFilterAnd, matchQuery)
		}
	}
}
