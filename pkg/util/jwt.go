package util

import (
	"time"

	jwt "github.com/dgrijalva/jwt-go"

	"github.com/EDDYCJY/go-gin-example/pkg/setting"
)

//jwtSecret 是我们设定的密钥
var jwtSecret = []byte(setting.JwtSecret)

//宣告JWT 結構
type Claims struct {
	Username           string `json:"username"`
	Password           string `json:"password"`
	jwt.StandardClaims        //通过 StandardClaims 生成标准的载体
}

//傳進player生成token
func GenerateToken(username, password string) (string, error) {
	nowTime := time.Now()
	//設定3小時過期
	expireTime := nowTime.Add(3 * time.Hour)
	//PAYLOAD
	claims := Claims{
		username,
		password,
		jwt.StandardClaims{
			ExpiresAt: expireTime.Unix(),
			Issuer:    "gin-blog",
		},
	}
	//宣告使用 HS256 與加入Payload 的聲明內容
	tokenClaims := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	//通过 HS256 算法生成 tokenClaims ,这就是我们的 HEADER 部分和 PAYLOAD。
	token, err := tokenClaims.SignedString(jwtSecret)

	return token, err
}

//通過 jwt-go 解析 jwt 就可達到用 jwt 驗證的效果
func ParseToken(token string) (*Claims, error) {
	tokenClaims, err := jwt.ParseWithClaims(token, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return jwtSecret, nil
	})

	if tokenClaims != nil {
		if claims, ok := tokenClaims.Claims.(*Claims); ok && tokenClaims.Valid {
			return claims, nil
		}
	}

	return nil, err
}
