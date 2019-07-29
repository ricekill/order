package common

import (
	"github.com/gin-gonic/gin"
)

type HandleFunc func(ctx *gin.Context) error

func Handle(f HandleFunc) gin.HandlerFunc {
	return func(c *gin.Context) {
		if err := f(c);err != nil {

		} else if c.Writer.Size() < 0 {
			
		}
	}
}