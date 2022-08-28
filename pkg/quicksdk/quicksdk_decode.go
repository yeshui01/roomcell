/*
 * @Author: mknight(tianyh)
 * @Mail: 824338670@qq.com
 * @Date: 2022-06-15 14:14:17
 * @LastEditTime: 2022-06-15 14:14:17
 * @Brief: quicksdk数据加密和解密
 */
package quicksdk

import (
	"regexp"
	"strconv"
)

func Decode(src string, cryptKey string) string {
	var des []byte
	// reg1 := regexp.MustCompile("(\\d+)")
	reg1 := regexp.MustCompile("\\d+")
	if reg1 == nil { //解释失败，返回nil
		return string(des)
	}
	//根据规则提取关键信息
	result := reg1.FindAllStringSubmatch(src, -1)
	if len(result) < 1 {
		return string(des)
	}
	var listN []int = make([]int, len(result))
	for i := 0; i < len(result); i++ {
		n, err := strconv.ParseInt(result[i][0], 10, 32)
		if err == nil {
			listN[i] = int(n)
		}
	}
	bKey := []byte(cryptKey)
	for i := 0; i < len(listN); i++ {
		c := rune(listN[i] - int(0xff&bKey[i%len(bKey)]))
		des = append(des, byte(c))
	}
	return string(des)
}
