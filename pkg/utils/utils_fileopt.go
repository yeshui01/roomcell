package utils

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"roomcell/pkg/loghlp"
	"strings"
)

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

// 检测文件是否存在
func CheckFileIsExist(filename string) bool {
	var exist = true
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		exist = false
	}
	return exist
}

// 保存内容到文件
func SaveJvDataToFile(jv interface{}, savefile string) {
	jsonData, errJv := json.MarshalIndent(jv, "", "\t")
	if errJv == nil {
		if CheckFileIsExist(savefile) {
			fileObj, fileErr := os.OpenFile(savefile, os.O_RDWR, 0666) //打开文件
			if fileErr != nil {
				return
			}
			fileObj.Truncate(0)
			fileObj.WriteString(string(jsonData))
			fileObj.Sync()
			fileObj.Close()
		} else {
			fileObj, fileErr := os.Create(fmt.Sprintf("%s", savefile)) //创建文件
			if fileErr == nil {
				fileObj.WriteString(string(jsonData))
				fileObj.Sync()
				fileObj.Close()
			}
		}
	} else {
		loghlp.Errorf("jv marshal error:%s", errJv)
	}
}

func ReadJvDataFromFile(savefile string, jv interface{}) error {
	if CheckFileIsExist(savefile) {
		fileObj, fileErr := os.Open(savefile) //打开文件
		if fileErr != nil {
			return errors.New("open file faile")
		}
		data, err := ioutil.ReadAll(fileObj)
		if err != nil {
			return err
		}
		errJv := json.Unmarshal(data, jv)
		return errJv
	}

	return errors.New("not find file")
}

// 提取文件名列表
func GetDirFiles(dirPath string) []string {
	var listData = make([]string, 0)
	_, err := os.Stat(dirPath)
	if err != nil {
		return listData
	}

	filepath.Walk(dirPath, func(filePath string, info os.FileInfo, err error) error {
		if info.IsDir() {
			return nil
		}
		lastIndex := strings.LastIndex(filePath, "\\")
		if lastIndex != -1 {
			listData = append(listData, filePath[lastIndex+1:])
			return nil
		}
		lastIndex = strings.LastIndex(filePath, "/")
		if lastIndex != -1 {
			listData = append(listData, filePath[lastIndex+1:])
			return nil
		}
		return nil
	})
	return listData
}
