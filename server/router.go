//Created by Goland
//@User: lenora
//@Date: 2021/2/5
//@Time: 2:40 下午
package server

import "github.com/gin-gonic/gin"

func (ds *defaultServer) routers() {
	exampleNone := ds.engine.Group("index")
	exampleNone.GET("", func(context *gin.Context) {
		ds.ResponseSuccess(context,nil)
		return
	})
}
