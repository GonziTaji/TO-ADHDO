package tags

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

var store Store = Store{}

type Controller struct {
}

func (Controller) GetListHandler(c *gin.Context) {
	limit := 10

	if limit_query := c.Query("limit"); limit_query != "" {
		l, err := strconv.Atoi(limit_query)

		if err == nil {
			limit = l
		} else {
			fmt.Printf("error parsing queryparam `limit`: %s\n", err.Error())
		}
	}

	tags, err := store.List(int8(limit), false)

	if err != nil {
		// TODO: handle error
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": err.Error(),
		})

		return
	}

	c.JSON(http.StatusOK, gin.H{
		"tags": tags,
	})

	return
}
