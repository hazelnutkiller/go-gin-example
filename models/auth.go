package models

import (
	u "github.com/EDDYCJY/go-gin-example/pkg/setting/e/util"
	"github.com/dgrijalva/jwt-go"
	"github.com/jinzhu/gorm"
	"golang.org/x/crypto/bcrypt"
)

//如何获取Token
//如何调用它呢，我们还要获取Token呢
//新增一个获取Token的API
//新增帳號結構
type Auth struct {
	ID       int    `gorm:"primary_key" json:"id"`
	Username string `json:"username"`
	Password string `json:"password"`
}

func CheckAuth(username, password string) bool {
	var auth Auth
	db.Select("id").Where(Auth{Username: username, Password: password}).First(&auth)
	if auth.ID > 0 {
		return true
	}
	return false
}

//認證玩家是否存在
func CheckUser(username string) bool {
	var auth Auth
	db.Select("id").Where(Auth{Username: username}).First(&auth)
	if auth.ID > 0 {
		return false
	}
	return true
}

//JWT 結構
type Token struct {
	UserId uint
	jwt.StandardClaims
}

//修改帳號結構
type Account struct {
	gorm.Model
	ID       int    `gorm:"primary_key" json:"id"`
	Username string `json:"username"`
	Password string `json:"Password"`
	Token    string `json:"token";sql:"-"`
}

//models\accounts.go 會包含驗證輸入的舊密碼是否正確，新密碼是否符合長度要求
//密碼修改
func Editor(oldPassword, newPassword, username string) map[string]interface{} {
	var auth Auth
	//var account Account
	//如果密碼小於8位數，回傳帳號格式錯誤
	if len(oldPassword) < 8 && len(newPassword) < 8 {
		return u.Message(false, "Password format error")
	}

	//帳號驗證
	// err := db.Table("id").Where(Auth{Username: username, Password: oldPassword}).First(username).Error
	// if err != nil {
	// 	if err == gorm.ErrRecordNotFound {
	// 		return u.Message(false, "username not found")
	// 	}
	// 	return u.Message(false, "Connection error")
	// }
	err := db.Select("id").Where("username = ?", username).First(&auth).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return u.Message(false, "username not found")
		}
		return u.Message(false, "Connection error")
	}

	//舊密碼驗證
	err = bcrypt.CompareHashAndPassword([]byte(oldPassword), []byte(auth.Password))
	if err != nil && err == bcrypt.ErrMismatchedHashAndPassword {
		return u.Message(false, "pwd error")
	}

	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	auth.Password = string(hashedPassword)

	//保存
	db.Save(auth)

	//清空 pwd 的數值以便回傳
	auth.Password = ""

	//修改完成
	response := u.Message(true, "Account has been editor")
	response["account"] = auth

	return response
	//回傳response給修改密碼的函數
}

//創建帳號存入DB
func AddUser(data map[string]interface{}) bool {
	db.Create(&Auth{
		Username: data["username"].(string),
		Password: data["password"].(string),
	})
	return true
}
