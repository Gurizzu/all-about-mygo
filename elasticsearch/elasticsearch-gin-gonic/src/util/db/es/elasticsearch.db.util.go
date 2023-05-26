package es

import "elasticsearch-gin-gonic/src/model"

func GetSkipAndLimit(requestPagination model.Request_Pagination) (skip int64, limit int64) {
	skip = requestPagination.Page
	if skip > 0 {
		skip--
	}
	skip *= requestPagination.Size

	limit = requestPagination.Size
	if limit == 0 {
		limit = 10
	}
	return
}
