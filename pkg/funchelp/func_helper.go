/* ====================================================================
 * Author           : tianyh(mknight)
 * Email            : 824338670@qq.com
 * Last modified    : 2022-07-25 11:33
 * Filename         : funchelp.go
 * Description      :
 * ====================================================================*/
package funchelp

import (
	"bufio"
	"bytes"
	"compress/zlib"
	"io"
	"os"
	"strings"
)

// string -> AaBb,默认下横线分割

func IsLowwerChar(c byte) bool {
	return c >= 'a' && c <= 'z'
}

func IsUpperChar(c byte) bool {
	return c >= 'A' && c <= 'Z'
}

func CharToUpper(c byte) byte {
	if IsLowwerChar(c) {
		return c - 32
	}
	return c
}

func CharToLower(c byte) byte {
	if IsUpperChar(c) {
		return c + 32
	}
	return c
}

func StringToAaBb(src string) string {
	if len(src) < 1 {
		panic("param error")
	}
	r := ""
	words := strings.Split(src, "_")
	for _, w := range words {
		if w == "id" {
			r = r + "ID"
		} else {
			s := []byte(w)
			if IsLowwerChar(s[0]) {
				s[0] = (CharToUpper(w[0]))
				r = r + string(s)
			} else {
				r = r + w
			}
		}
	}
	return r
}

func StringToaaBb(src string) string {
	if len(src) < 1 {
		panic("param error")
	}
	r := ""
	words := strings.Split(src, "_")
	for i, w := range words {
		if i == 0 {
			if w == "ID" || w == "iD" || w == "Id" {
				r = r + "id"
				continue
			}
			s := []byte(w)
			if IsUpperChar(s[0]) {
				s[0] = (CharToLower(w[0]))
				r = r + string(s)
			} else {
				r = r + w
			}
		} else {
			r = r + StringToAaBb(w)
		}
	}
	return r
}

// 读取文件内容为行列表
func ReadFileForLine(fileName string) [][]byte {
	f, err := os.Open(fileName)
	if err != nil {
		return nil
	}
	lineList := make([][]byte, 0)
	fread := bufio.NewReader(f)
	for {
		line, _, err := fread.ReadLine()
		if err == io.EOF {
			break
		}
		copyLine := make([]byte, len(line))
		copy(copyLine, line)
		lineList = append(lineList, copyLine)
	}
	return lineList
}
func ReadFileForLineStr(fileName string) []string {
	f, err := os.Open(fileName)
	if err != nil {
		return make([]string, 0)
	}
	lineList := make([]string, 0)
	fread := bufio.NewReader(f)
	for {
		line, _, err := fread.ReadLine()
		if err == io.EOF {
			break
		}
		copyLine := make([]byte, len(line))
		copy(copyLine, line)
		lineList = append(lineList, string(copyLine))
	}
	return lineList
}
func CheckFileIsExist(filename string) bool {
	var exist = true
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		exist = false
	}
	return exist
}

//进行zlib压缩
func ZlibCompress(srcData []byte) []byte {
	var in bytes.Buffer
	w := zlib.NewWriter(&in)
	w.Write(srcData)
	w.Close()
	return in.Bytes()
}

//进行zlib解压缩
func ZlibUnCompress(compressedData []byte) []byte {
	b := bytes.NewReader(compressedData)
	var out bytes.Buffer
	r, _ := zlib.NewReader(b)
	io.Copy(&out, r)
	return out.Bytes()
}
