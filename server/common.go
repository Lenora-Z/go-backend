package server

import "github.com/gin-gonic/gin"

func (ds *defaultServer) getLanguage(c *gin.Context) string {
	return c.Request.Header.Get("Lang")
}
