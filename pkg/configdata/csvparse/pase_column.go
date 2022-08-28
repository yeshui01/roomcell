package csvparse

import (
	"bufio"
	"io"
	"os"
	"strings"

	"github.com/spf13/cast"
)

func ParseColumnInt(strColumn string) int32 {
	return cast.ToInt32(strColumn)
}
func ParseColumnString(strColumn string) string {
	return strColumn
}
func ParseColumnArrayInt(strColumn string) []int32 {
	strList := strings.Split(strColumn, "|")
	intList := make([]int32, len(strList))
	for i, s := range strList {
		intList[i] = cast.ToInt32(s)
	}
	return intList
}
func ParseColumnArrayString(strColumn string) []string {
	strList := strings.Split(strColumn, "|")
	return strList
}

func ParseCsvFileLine(fileName string) []string {
	f, err := os.Open(fileName)
	if err != nil {
		return make([]string, 0)
	}
	lineList := make([]string, 0)
	fread := bufio.NewReader(f)
	lineIdx := 0
	for {
		line, _, err := fread.ReadLine()
		if err == io.EOF {
			break
		}
		if lineIdx < 4 { // 忽略head的4行,head的4行为字段信息
			lineIdx++
			continue
		}
		copyLine := make([]byte, len(line))
		copy(copyLine, line)
		lineList = append(lineList, string(copyLine))
		lineIdx++
	}
	return lineList
}
