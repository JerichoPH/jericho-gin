package tools

import (
	"jericho-go/wrongs"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

var jwtSecret = []byte("jwt_key") // 加密密钥

// Claims Jwt 表单
type Claims struct {
	Id       uint64 `json:"id"`
	Account  string `json:"account"`
	Nickname string `json:"nickname"`
	Uuid     string `json:"uuid"`
	jwt.StandardClaims
}

// GenerateJwt 生成Jwt标准过期时间
func GenerateJwt(
	id uint64,
	account string,
	nickname string,
	uuid string,
) (string, error) {
	return generateJwtWithExpire(
		id,
		account,
		nickname,
		uuid,
		24*180*time.Hour,
	)
}

// GenerateJwtLongTime 生成Jwt长时间过期时间
func GenerateJwtLongTime(
	id uint64,
	account string,
	nickname string,
	uuid string,
) (string, error) {
	return generateJwtWithExpire(
		id,
		account,
		nickname,
		uuid,
		292*24*365*time.Hour,
	)
}

// 通过过期时间生成jwt
func generateJwtWithExpire(
	id uint64,
	account string,
	nickname string,
	uuid string,
	expireOffsetTime time.Duration,
) (string, error) {
	// 设置token有效时间
	nowTime := time.Now()
	expireTime := nowTime.Add(time.Duration(expireOffsetTime))

	claims := Claims{
		Id:       id,
		Account:  account,
		Nickname: nickname,
		Uuid:     uuid,
		StandardClaims: jwt.StandardClaims{
			// 过期时间
			ExpiresAt: expireTime.Unix(),
			// 指定token发行人
			Issuer: "jericho-go",
		},
	}

	tokenClaims := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	// 该方法内部生成签名字符串，再用于获取完整、已签名的token
	token, err := tokenClaims.SignedString(jwtSecret)
	return token, err
}

// ParseJwt 根据传入的token值获取到Claims对象信息，（进而获取其中的用户名和密码）
func ParseJwt(token string) (*Claims, error) {
	// 用于解析鉴权的声明，方法内部主要是具体的解码和校验的过程，最终返回*Token
	tokenClaims, err := jwt.ParseWithClaims(token, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return jwtSecret, nil
	})

	if tokenClaims != nil {
		// 从tokenClaims中获取到Claims对象，并使用断言，将该对象转换为我们自己定义的Claims
		// 要传入指针，项目中结构体都是用指针传递，节省空间。
		if claims, ok := tokenClaims.Claims.(*Claims); ok && tokenClaims.Valid {
			return claims, nil
		}
	}
	return nil, err
}

// GetJwtFromHeader 从header中获取jwt
func GetJwtFromHeader(ctx *gin.Context) string {
	tokens := ctx.Request.Header["Authorization"]

	if len(tokens) == 0 {
		wrongs.ThrowUnAuth("令牌不存在")
	}
	return tokens[0]
}
