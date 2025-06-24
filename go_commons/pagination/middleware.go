package pagination

import (
	"github.com/gin-gonic/gin"
	"strconv"
)

func Middleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set(Limit, ParsePerPage(c.Query(PerPage)))
		c.Set(Page, ParsePage(c.Query(Page)))

		c.Next()
	}
}

func ParsePerPage(value string) int64 {
	perPage, err := strconv.ParseInt(value, 10, 0)
	if err != nil || perPage > 100 {
		perPage = 20
	}

	return perPage
}

func ParsePage(value string) int64 {
	page, err := strconv.ParseInt(value, 10, 0)
	if err != nil {
		page = 1
	}

	return page
}
