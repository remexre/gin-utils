package ginUtils

import (
	"database/sql"

	"github.com/gin-gonic/gin"
)

// A middleware that adds the DB to the context.
func SetDB(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set("db", db)
		c.Next()
	}
}
