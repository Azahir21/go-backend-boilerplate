package http

import "github.com/gin-gonic/gin"

func GetRequest[T any](c *gin.Context) *T {
	v, ok := c.Get("req")
	if !ok {
		return nil
	}
	if cast, ok := v.(*T); ok {
		return cast
	}
	if cast, ok := v.(T); ok {
		return &cast
	}
	return nil
}
