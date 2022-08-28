package crossdef

import (
	"errors"
	"time"

	"github.com/dgrijalva/jwt-go"
)

// JWT 签名结构
type JWT struct {
	SigningKey []byte
}

// 一些常量
var (
	ErrTokenExpired     error  = errors.New("token is expired")
	ErrTokenNotValidYet error  = errors.New("token not active yet")
	ErrTokenMalformed   error  = errors.New("that's not even a token")
	ErrTokenInvalid     error  = errors.New("couldn't handle this token")
	SignKey             string = "123456"
)

// 载荷，可以加一些自己需要的信息
type CustomClaims struct {
	Account       string `json:"account"`
	LastLoginTime int64  `json:"last_login_time"`
	CchId         string `json:"cch_id"`
	UserID        int64  `json:"user_id"`
	DataZone      int32  `json:"data_zone"`
	jwt.StandardClaims
}

// 新建一个jwt实例
func NewJWT() *JWT {
	return &JWT{
		[]byte(GetSignKey()),
	}
}

// 获取signKey
func GetSignKey() string {
	return SignKey
}

// 这是SignKey
func (j *JWT) SetSignKey(key string) {
	j.SigningKey = []byte(key)
}

// CreateToken 生成一个token
func (j *JWT) CreateToken(claims CustomClaims) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(j.SigningKey)
}

// 解析Tokne
func (j *JWT) ParseToken(tokenString string) (*CustomClaims, error) {
	customRet := &CustomClaims{}
	token, err := jwt.ParseWithClaims(tokenString, customRet, func(token *jwt.Token) (interface{}, error) {
		return j.SigningKey, nil
	})
	if err != nil {
		if ve, ok := err.(*jwt.ValidationError); ok {
			if ve.Errors&jwt.ValidationErrorMalformed != 0 {
				return nil, ErrTokenMalformed
			} else if ve.Errors&jwt.ValidationErrorExpired != 0 {
				// Token is expired
				if token != nil {
					return token.Claims.(*CustomClaims), ErrTokenExpired
				}
				return customRet, ErrTokenExpired
			} else if ve.Errors&jwt.ValidationErrorNotValidYet != 0 {
				return customRet, ErrTokenNotValidYet
			} else {
				return customRet, ErrTokenInvalid
			}
		}
	}
	if claims, ok := token.Claims.(*CustomClaims); ok && token.Valid {
		return claims, nil
	}
	return customRet, ErrTokenInvalid
}

// 更新token
func (j *JWT) RefreshToken(tokenString string) (string, error) {
	jwt.TimeFunc = func() time.Time {
		return time.Unix(0, 0)
	}
	token, err := jwt.ParseWithClaims(tokenString, &CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		return j.SigningKey, nil
	})
	if err != nil {
		return "", err
	}
	if claims, ok := token.Claims.(*CustomClaims); ok && token.Valid {
		jwt.TimeFunc = time.Now
		claims.StandardClaims.ExpiresAt = time.Now().Add(1 * time.Hour).Unix()
		return j.CreateToken(*claims)
	}
	return "", ErrTokenInvalid
}
func TokenAuthClaims(token string, signKey string) (bool, *CustomClaims) {
	j := NewJWT()
	j.SetSignKey(signKey)
	cus, parse_err := j.ParseToken(token)
	if parse_err != nil {
		return false, cus
	}
	return true, cus
}
