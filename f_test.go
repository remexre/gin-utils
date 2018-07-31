package ginUtils_test

import (
	"database/sql"

	"github.com/gin-gonic/gin"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	. "github.com/remexre/gin-utils"
)

var _ = Describe("F", func() {
	Context("Signature Checking", func() {
		It("Should panic if passed a non-function", func() {
			Expect(func() { F("not a function") }).To(Panic())
		})

		It("Should panic if passed a function with the wrong arity", func() {
			Expect(func() {
				F(func(ctx *gin.Context) (int, interface{}) { return 0, nil })
			}).To(Panic())
		})

		It("Should panic if passed a function that doesn't take a context", func() {
			Expect(func() {
				F(func(ctx, db *sql.DB) (int, interface{}) { return 0, nil })
			}).To(Panic())
		})

		It("Should panic if passed a variadic function", func() {
			Expect(func() {
				F(func(ctx *gin.Context, db ...*sql.DB) (int, interface{}) { return 0, nil })
			}).To(Panic())
		})

		It("Should panic if passed a function with the wrong number of return values", func() {
			Expect(func() {
				F(func(ctx *gin.Context, db ...*sql.DB) int { return 0 })
			}).To(Panic())
		})

		It("Should panic if passed a function whose first return value isn't an int", func() {
			Expect(func() {
				F(func(ctx *gin.Context, db ...*sql.DB) (string, interface{}) { return "", nil })
			}).To(Panic())
		})

		It("Should not panic if passed a function with the appropriate signature", func() {
			// Try multiple DB types.
			Expect(func() {
				F(func(ctx *gin.Context, db *sql.DB) (int, interface{}) { return 0, nil })
			}).NotTo(Panic())

			Expect(func() {
				F(func(ctx *gin.Context, db gin.H) (int, interface{}) { return 0, nil })
			}).NotTo(Panic())

			Expect(func() {
				F(func(ctx *gin.Context, db interface{}) (int, interface{}) { return 0, nil })
			}).NotTo(Panic())
		})

		It("Should not panic if passed a function with a more specific return type", func() {
			Expect(func() {
				F(func(ctx *gin.Context, db *sql.DB) (int, gin.H) { return 0, nil })
			}).NotTo(Panic())

			Expect(func() {
				F(func(ctx *gin.Context, db *sql.DB) (int, string) { return 0, "" })
			}).NotTo(Panic())
		})
	})
})
