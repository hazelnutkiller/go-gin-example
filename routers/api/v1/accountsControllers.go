package v1

import (
	"encoding/json"

	"net/http"

	"github.com/EDDYCJY/go-gin-example/pkg/setting/e"
	"github.com/gin-gonic/gin"
)

//修改密碼的函數
var EditorAccount = func(c *gin.Context) {

	account := &models.Account{}                   //引用 accounts.go
	err := json.NewDecoder(r.Body).Decode(account) //解析傳入的 json 請求

	//如果輸入的請求錯誤
	if err != nil {
		code = e.INVALID_REQUEST
		return
	} else {
		if code != e.SUCCESS {
			c.JSON(http.StatusUnauthorized, gin.H{
				"code": code,
				"msg":  e.GetMsg(code),
				"data": data,
			})

		}

		resp := account.Editor() //呼叫 accounts.go 透過傳入的資料刪除帳號
		u.Respond(w, resp)
	}
}
