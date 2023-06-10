package v1

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"

	"github.com/EDDYCJY/go-gin-example/models"
	"github.com/EDDYCJY/go-gin-example/pkg/setting/e"
	"github.com/astaxie/beego/validation"
	"github.com/gin-gonic/gin"
)

type Account struct {
	Id       int64  `json:"id" form:"id"`
	Username string `valid:"Required; MaxSize(50);"`
	Password string `valid:"Required; MaxSize(50)"`
}

func AddUser(c *gin.Context) {
	urlstr := "http://127.0.0.1:8000/api/v1/adduser"
	username := c.Query("username")
	password := c.Query("password")
	values := url.Values{}
	valid := validation.Validation{}
	valid.Required(username, "username").Message("帳號不能為空")
	valid.Required(password, "password").Message("密碼不能為空")
	values.Set("username", username)
	values.Set("password", password)

	code := e.ERROR_PLAYER_EXISTS

	if !valid.HasErrors() {
		if models.CheckUser(username) {
			data := make(map[string]interface{})
			data["username"] = username
			data["password"] = password

			models.AddUser(data)
			code = e.SUCCESS
		} else {
			for _, err := range valid.Errors {
				log.Printf("err.key: %s, err.message: %s", err.Key, err.Message)
			}
		}

		r, err := http.PostForm(urlstr, values)
		if err != nil {
			log.Fatal(err)
		}
		defer r.Body.Close()
		//读取整个响应体
		body, _ := ioutil.ReadAll(r.Body)
		var data Account
		json.Unmarshal([]byte(body), &data)

		c.JSON(http.StatusOK, gin.H{
			"code": code,
			"msg":  e.GetMsg(code),
			"data": data,
		})
	}
}

type EditorUesr struct {
	Id          int64  `json:"id" form:"id"`
	Username    string `valid:"Required; MaxSize(50);"`
	OldPassword string `valid:"Required; MaxSize(50)"`
	NewPassword string `valid:"Required; MaxSize(50)"`
}

//修改密碼的函數
func EditorAccount(c *gin.Context) {
	//var editoruesr EditorUesr
	urlstr := "http://127.0.0.1:8000/api/v1/user/editor"
	values := url.Values{}
	//valid := validation.Validation{}
	Username := c.Query("username")
	OldPassword := c.Query("oldPassword")
	NewPassword := c.Query("newPassword")
	//引用auth結構
	//account := &models.Auth{}
	data := make(map[string]interface{})
	data["username"] = Username
	data["oldPassword"] = OldPassword
	data["newPassword"] = NewPassword

	models.Editor(Username, OldPassword, NewPassword)

	code := e.SUCCESS

	r, err := http.PostForm(urlstr, values)
	if err != nil {
		log.Fatal(err)
	}
	defer r.Body.Close()
	//读取整个响应体
	body, _ := ioutil.ReadAll(r.Body)

	json.Unmarshal([]byte(body), &data)

	c.JSON(http.StatusOK, gin.H{
		"code": code,
		"msg":  e.GetMsg(code),
		"data": data,
	})

}
