/*
 * @Author: mknight(tianyh)
 * @Mail: 824338670@qq.com
 * @Date: 2022-10-08 15:15:19
 * @LastEditTime: 2022-10-10 13:59:28
 * @FilePath: \roomcell\app\account\accrouter\account_token.go
 */
package accrouter

import (
	"roomcell/pkg/crossdef"
	"time"

	"github.com/dgrijalva/jwt-go"
)

func genJwtToken(userID int64, userName string, nickName string, dataZone int32) (string, error) {
	j := &crossdef.JWT{
		[]byte(crossdef.SignKey),
	}
	tokenTime := int64(3600 * 24) // 默认24小时
	claims := crossdef.CustomClaims{
		Account:       userName,
		Nickname:      nickName,
		LastLoginTime: 0,
		CchId:         "",
		DataZone:      dataZone,
		UserID:        userID,
		StandardClaims: jwt.StandardClaims{
			NotBefore: time.Now().Unix() - 1000,      // 签名生效时间
			ExpiresAt: time.Now().Unix() + tokenTime, // 过期时间 一小时
			Issuer:    crossdef.SignKey,              //签名的发行者
		},
	}

	token, err := j.CreateToken(claims)
	return token, err
}
