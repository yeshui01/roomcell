/*
 * @Author: mknight(tianyh)
 * @Mail: 824338670@qq.com
 * @Date: 2022-06-15 14:14:17
 * @LastEditTime: 2022-06-15 14:14:17
 * @Brief: quicksdk数据加密和解密
 */
package quicksdk

import "fmt"

func Encode(src string, cryptKey string) string {
	var des string
	bKey := []byte(cryptKey)
	for i := 0; i < len(src); i++ {
		n := (0xff & src[i]) + (0xff & bKey[i%len(bKey)])
		groupStr := fmt.Sprintf("@%d", n)
		des = des + groupStr
	}
	return string(des)
}
