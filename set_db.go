package ginUtils

import (
	"github.com/gin-gonic/gin"
)

// SetDB is a middleware that adds the DB to the context.
func SetDB(db interface{}) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set("db", db)
		c.Next()
	}
}
