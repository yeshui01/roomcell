/* ====================================================================
 * Author           : tianyh(mknight)
 * Email            : 824338670@qq.com
 * Last modified    : 2022-07-25 11:33
 * Filename         : funchelp.go
 * Description      : 读取excel配置文件,导出csv文件
 * ====================================================================*/
package csvtools

import (
	"fmt"
	"os"
	"roomcell/pkg/funchelp"

	"github.com/360EntSecGroup-Skylar/excelize"
	"github.com/sirupsen/logrus"
)

func fetchFileNameWithoutExtension(excelFile string) string {
	strLen := len(excelFile)
	idxEnd := strLen - 1
	for {
		if idxEnd <= 0 {
			break
		}
		if excelFile[idxEnd] == '.' {
			break
		} else {
			idxEnd--
		}
	}
	idxBegin := idxEnd
	for {
		if idxBegin <= 0 {
			break
		}
		if excelFile[idxBegin] == '\\' || excelFile[idxBegin] == '/' {
			break
		}
		idxBegin--
	}
	return string([]byte(excelFile)[idxBegin+1 : idxEnd])
}

func ReadExcelToCSV(excelFile string, csvSavePath string) {
	f, err := excelize.OpenFile(excelFile)
	if err != nil {
		logrus.Errorf("open file[%s]error:%s", excelFile, err.Error())
		return
	}

	var fileObj *os.File
	var fileErr error
	fileName := fmt.Sprintf("%s/%s.csv", csvSavePath, fetchFileNameWithoutExtension(excelFile))
	if funchelp.CheckFileIsExist(fileName) {
		fileObj, fileErr = os.OpenFile(fileName, os.O_RDWR, 0666) //打开文件
		if fileErr != nil {
			return
		}
	} else {
		fileObj, fileErr = os.Create(fileName) //创建文件
		if fileErr != nil {
			return
		}
	}
	defer fileObj.Close()
	fileObj.Truncate(0)
	rows := f.GetRows("Sheet1") // 默认读取第一个表单
	if len(rows) < 4 {
		return
	}
	resultRows := make([][]string, len(rows))
	for r, oneRow := range rows {
		newRow := make([]string, 0)
		// head 部分检测
		for c, oneColumn := range oneRow {
			if len(rows[0][c]) == 0 || len(rows[1][c]) == 0 || len(rows[2][c]) == 0 || len(rows[3][c]) == 0 {
				continue
			}
			if len(oneColumn) == 0 {
				newRow = append(newRow, "0")
			} else {
				newRow = append(newRow, oneColumn)
			}
		}
		resultRows[r] = newRow
	}
	totalLine := len(resultRows)
	for i, oneRow := range resultRows {
		var lineStr string
		for i, cellStr := range oneRow {
			if i == 0 {
				lineStr = cellStr
			} else {
				lineStr = lineStr + "," + cellStr
			}
		}
		if i == totalLine-1 {
			fileObj.WriteString(lineStr)
		} else {
			fileObj.WriteString(lineStr + "\n")
		}
	}
	fileObj.Sync()
}
