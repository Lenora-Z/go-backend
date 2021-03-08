//Created by Goland
//@User: lenora
//@Date: 2021/3/8
//@Time: 11:04 上午
package server

import (
	"git.dataqin.com/lenora/go-backend/utils"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
	"time"
)

func (srv *defaultServer) CreateToken(uId uint16, mobile, name string) (string, error) {
	at := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id":   uId,
		"exp":       time.Now().Add(utils.TokenExpire).Unix(),
		"mobile":    mobile,
		"user_name": name,
		"uu_id":     "",
	})

	token, err := at.SignedString([]byte(srv.conf.JWTSecret))
	if err != nil {
		return "", err
	}
	return token, nil
}

func (srv *defaultServer) CheckToken() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		header := ctx.GetHeader("Authorization")
		g := strings.Split(header, " ")
		if g[0] != "Bearer" || g[1] == "" {
			ctx.JSON(http.StatusUnauthorized, ApiResponse{
				Code: 2002,
			})
			ctx.Abort()
			return
		}
		ret, err := jwt.Parse(g[1], func(token *jwt.Token) (interface{}, error) {
			return []byte(srv.conf.JWTSecret), nil
		})
		if err != nil {
			ctx.JSON(http.StatusUnauthorized, ApiResponse{
				Code: 2002,
			})
			ctx.Abort()
			return
		}

		cla := ret.Claims.(jwt.MapClaims)
		expire := cla["exp"].(float64)

		if expire < float64(time.Now().Unix()) {
			ctx.JSON(http.StatusOK, ApiResponse{
				Code:    2001,
				Message: "user licence expired",
			})
			ctx.Abort()
			return
		}
		claims := CustomClaims{
			UserID:   uint64(cla["user_id"].(float64)),
			Name:     cla["user_name"].(string),
			UserName: cla["user_name"].(string),
		}

		ctx.Set("claims", &claims)
		ctx.Next()
	}
}
