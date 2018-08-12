package ginUtils

import (
	"fmt"
	"log"
	"net/http"
	"reflect"

	"github.com/gin-gonic/gin"
)

type redirect struct{ url string }

// Redirect returns an object that when returned from an F handler, represents a redirect.
func Redirect(url string) interface{} {
	return redirect{url}
}

// F makes a wrapper for a handler-like function that adds some niceties. The wrapped function
// should return an HTTP status and a response object. If the response object is nil, no response
// will be written. If the response object is returned by Redirect, a 307 Temporary Redirect will
// be issued.
//
// The handler's type should be of the form func(*gin.Context, T) (int, interface{}), where T is
// the same type that was passed to the SetDB middleware.
//
// Requires that the SetDB middleware be mounted.
func F(handler interface{}) gin.HandlerFunc {
	handlerValue := reflect.ValueOf(handler)
	handlerType := handlerValue.Type()
	if handlerType.Kind() != reflect.Func {
		panic(fmt.Sprintf("Handler %v is not a function", handler))
	}
	if handlerType.NumIn() != 2 {
		panic(fmt.Sprintf("Handler %v takes the wrong number of arguments", handler))
	}
	if handlerType.In(0) != reflect.TypeOf((*gin.Context)(nil)) {
		panic(fmt.Sprintf("Handler %v's first argument should be *gin.Context", handler))
	}
	if handlerType.IsVariadic() {
		panic(fmt.Sprintf("Handler %v should not be variadic", handler))
	}
	if handlerType.NumOut() != 2 {
		panic(fmt.Sprintf("Handler %v returns the wrong number of results", handler))
	}
	if handlerType.Out(0) != reflect.TypeOf(int(0)) {
		panic(fmt.Sprintf("Handler %v's first return value should be int, not %v", handler,
			handlerType.Out(0)))
	}

	handlerClo := func(c *gin.Context, db interface{}) (int, interface{}) {
		// We've already checked the type, so this is safe.
		retVals := handlerValue.Call([]reflect.Value{
			reflect.ValueOf(c),
			reflect.ValueOf(db),
		})
		respCode := int(retVals[0].Int())
		respData := retVals[1].Interface()
		return respCode, respData
	}

	return func(c *gin.Context) {
		db, has := c.Get("db")
		if !has {
			panic("db not present; was the SetDB middleware not used?")
		}

		status, body := handlerClo(c, db)
		if body != nil {
			respondWith(c, status, body)
		}
	}
}

func intOr(a, b int) int {
	if a == 0 {
		return b
	}
	return a
}

func respondWith(c *gin.Context, status int, body interface{}) {
	if err, ok := body.(error); ok {
		log.Println(err)
		c.JSON(intOr(status, http.StatusInternalServerError), gin.H{
			"status": "error",
			"err":    err.Error(),
		})
	} else if redir, ok := body.(redirect); ok {
		c.Redirect(intOr(status, http.StatusTemporaryRedirect), redir.url)
	} else {
		c.JSON(intOr(status, http.StatusOK), body)
	}
}
