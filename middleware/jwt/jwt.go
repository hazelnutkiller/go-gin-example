package jwt

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/EDDYCJY/go-gin-example/pkg/setting/e"
	"github.com/EDDYCJY/go-gin-example/pkg/util"
)

// 定义一个JWTAuth的中间件
func JWT() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 通过http header中的token解析来认证
		var code int
		var data interface{}

		code = e.SUCCESS
		token := c.Query("token")
		if token == "" {
			code = e.INVALID_PARAMS
		} else {
			claims, err := util.ParseToken(token)
			if err != nil {
				code = e.ERROR_AUTH_CHECK_TOKEN_FAIL
			} else if time.Now().Unix() > claims.ExpiresAt {
				code = e.ERROR_AUTH_CHECK_TOKEN_TIMEOUT
			}
		}

		if code != e.SUCCESS {
			c.JSON(http.StatusUnauthorized, gin.H{
				"code": code,
				"msg":  e.GetMsg(code),
				"data": data,
			})

			c.Abort()
			return
		}

		c.Next()
	}
}
