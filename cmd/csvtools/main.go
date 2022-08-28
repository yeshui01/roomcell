package main

import (
	"fmt"
	"os"
	"roomcell/app/csvtools"
	"roomcell/pkg/utils"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/flopp/go-findfont"
)

type CSVToolsConfig struct {
	PkgPath    string `json:"pkg_path"`
	ExcelPath  string `json:"excel_path"`
	ExcelFile  string `json:"excel_file"`
	CsvPath    string `json:"csv_path"`
	ModuleName string `json:"module_name"`
}

// 引用中文
func init() {
	fontPaths := findfont.List()
	for i, fpath := range fontPaths {
		// 楷体: simkai
		// 黑体: simhei
		if strings.Contains(fpath, "simkai.ttf") {
			os.Setenv("FYNE_FONT", fpath)
			fmt.Println(i)
			break
		}
	}
}

func main() {
	// testExcelFile := "E:\\tyh_work_card\\svn_fp\\design\\excel\\union_boss.xlsx"
	// pkgPath := "E:\\tyh_work_card\\server_project\roomcell\\pkg"
	// csvPath := "E:\tyh_work_card\\server_project\\roomcell\\bin\\csv"
	// csvtools.ReadExcelToCSV(testExcelFile, csvPath)
	// csvtools.GenerageCSVDefCode("roomcell", pkgPath, testExcelFile)
	// csvtools.CSVModuleJoinToConfigMgr(testExcelFile, pkgPath)
	// csvtools.ExportConfigMgrInterface(testExcelFile, pkgPath)
	// configdata.InitConfigData(csvPath)

	configPath := "./"
	configSetting := &CSVToolsConfig{}
	// 读取配置
	utils.ReadJvDataFromFile(configPath+"/csv_tools.json", configSetting)
	myApp := app.New()
	myWindow := myApp.NewWindow("csv tools")
	myApp.Settings().SetTheme(theme.DarkTheme())
	pkgPathEntry := widget.NewEntry()
	if len(configSetting.PkgPath) > 0 {
		pkgPathEntry.SetPlaceHolder(configSetting.PkgPath)
		pkgPathEntry.SetText(configSetting.PkgPath)
	} else {
		pkgPathEntry.SetPlaceHolder("Enter module pkg path...")
	}

	excelPathEntry := widget.NewEntry()
	if len(configSetting.ExcelPath) > 0 {
		excelPathEntry.SetPlaceHolder(configSetting.ExcelPath)
		excelPathEntry.SetText(configSetting.ExcelPath)
	} else {
		excelPathEntry.SetPlaceHolder("Enter excel dir path...")
	}

	csvPathEntry := widget.NewEntry()
	if len(configSetting.CsvPath) > 0 {
		csvPathEntry.SetPlaceHolder(configSetting.CsvPath)
		csvPathEntry.SetText(configSetting.CsvPath)
	} else {
		csvPathEntry.SetPlaceHolder("Enter module csv path...")
	}

	input := widget.NewEntry()
	if len(configSetting.ExcelFile) > 0 {
		input.SetPlaceHolder(configSetting.ExcelFile)
		input.SetText(configSetting.ExcelFile)
	} else {
		input.SetPlaceHolder("Enter config excel...")
	}

	moduleNameEntry := widget.NewEntry()
	if len(configSetting.ModuleName) > 0 {
		moduleNameEntry.SetPlaceHolder(configSetting.ModuleName)
		moduleNameEntry.SetText(configSetting.ModuleName)
	} else {
		moduleNameEntry.SetPlaceHolder("Enter go module name...")
	}

	infoLabel := widget.NewLabel("")
	form := &widget.Form{
		Items: []*widget.FormItem{ // we can specify items in the constructor
			{Text: "pkg路径  :", Widget: pkgPathEntry},
			{Text: "excel路径:", Widget: excelPathEntry},
			{Text: "csv路径  :", Widget: csvPathEntry},
			{Text: "模块名字:", Widget: moduleNameEntry},
			{Text: "excel文件:", Widget: input},
		},
		OnSubmit: func() { // optional, handle form submission
			testExcelFile := input.Text
			pkgPath := pkgPathEntry.Text
			csvPath := csvPathEntry.Text
			excelPath := excelPathEntry.Text
			moduleName := moduleNameEntry.Text

			if len(excelPath) > 0 && len(moduleName) > 0 && len(csvPath) > 0 && len(pkgPath) > 0 && len(testExcelFile) > 0 {
				// utils.SaveJvDataToFile(&CSVToolsConfig{
				// 	ExcelFile:  testExcelFile,
				// 	PkgPath:    pkgPath,
				// 	CsvPath:    csvPath,
				// 	ExcelPath:  excelPath,
				// 	ModuleName: moduleName,
				// }, "./csv_tools.json")
				excelFileFullPath := fmt.Sprintf("%s/%s", configSetting.ExcelPath, testExcelFile)
				csvtools.ReadExcelToCSV(excelFileFullPath, csvPath)
				csvtools.GenerageCSVDefCode(moduleName, pkgPath, excelFileFullPath)
				csvtools.CSVModuleJoinToConfigMgr(excelFileFullPath, pkgPath)
				csvtools.ExportConfigMgrInterface(excelFileFullPath, pkgPath)
				//configdata.InitConfigData(csvPath)
				infoLabel.SetText("生成成功!")
			} else {
				infoLabel.SetText("配置设置错误!!!!")
			}
		},
		SubmitText: "生成",
		// OnCancel: func() {
		// 	testExcelFile := input.Text
		// 	pkgPath := pkgPathEntry.Text
		// 	csvPath := csvPathEntry.Text
		// 	excelPath := excelPathEntry.Text
		// 	moduleName := moduleNameEntry.Text
		// 	if len(excelPath) > 0 && len(moduleName) > 0 && len(csvPath) > 0 && len(pkgPath) > 0 && len(testExcelFile) > 0 {
		// 		utils.SaveJvDataToFile(&CSVToolsConfig{
		// 			ExcelFile:  testExcelFile,
		// 			PkgPath:    pkgPath,
		// 			CsvPath:    csvPath,
		// 			ExcelPath:  excelPath,
		// 			ModuleName: moduleName,
		// 		}, "./csv_tools.json")
		// 		infoLabel.SetText("配置保存成功!")
		// 	} else {
		// 		infoLabel.SetText("配置设置错误!!!!")
		// 	}
		// },
		// CancelText: "保存配置",
	}
	listData := utils.GetDirFiles(configSetting.ExcelPath)
	saveSettingBtn := widget.NewButton("   保存配置  ", func() {
		testExcelFile := input.Text
		pkgPath := pkgPathEntry.Text
		csvPath := csvPathEntry.Text
		excelPath := excelPathEntry.Text
		moduleName := moduleNameEntry.Text
		if len(excelPath) > 0 && len(moduleName) > 0 && len(csvPath) > 0 && len(pkgPath) > 0 && len(testExcelFile) > 0 {
			utils.SaveJvDataToFile(&CSVToolsConfig{
				ExcelFile:  testExcelFile,
				PkgPath:    pkgPath,
				CsvPath:    csvPath,
				ExcelPath:  excelPath,
				ModuleName: moduleName,
			}, "./csv_tools.json")
			infoLabel.SetText("配置保存成功!")
			// listData = utils.GetDirFiles(configSetting.ExcelPath)
			configSetting.CsvPath = csvPath
			configSetting.ExcelFile = testExcelFile
			configSetting.ExcelPath = excelPath
			configSetting.ModuleName = moduleName
			configSetting.PkgPath = pkgPath
		} else {
			infoLabel.SetText("配置设置错误!!!!")
		}
	})
	fileNameSearch := widget.NewEntry()
	fileNameSearch.SetPlaceHolder("输入要搜索的文件名")

	// 文件列表
	fileList := widget.NewList(
		func() int {
			return len(listData)
		},
		func() fyne.CanvasObject {
			return widget.NewLabel("template")
		},
		func(i widget.ListItemID, o fyne.CanvasObject) {
			o.(*widget.Label).SetText(listData[i])
		})
	searchBtn := widget.NewButton("   搜索文件     ", func() {
		if len(fileNameSearch.Text) > 0 {
			listDataMatchs := make([]string, 0)
			for _, v := range listData {
				if strings.Index(v, fileNameSearch.Text) != -1 {
					listDataMatchs = append(listDataMatchs, v)
				}
			}
			listData, listDataMatchs = listDataMatchs, listData
			fileList.Refresh()
		}
	})
	cancelSearchBtn := widget.NewButton("   取消搜索     ", func() {
		listData = utils.GetDirFiles(configSetting.ExcelPath)
		fileList.Refresh()
	})
	refreshBtn := widget.NewButton("    刷新   ", func() {
		listData = utils.GetDirFiles(configSetting.ExcelPath)
		fileList.Refresh()
	})
	searchContent := container.NewVBox(saveSettingBtn, fileNameSearch, searchBtn, cancelSearchBtn, refreshBtn)

	// // 子窗口
	// fileListWindow := myApp.NewWindow("fileList")
	// fileListWindow.SetOnClosed(func() {
	// })
	// leftPos := widget.NewLabel("")
	// topPos := widget.NewLabel("配置列表")
	// rightBtn := widget.NewButton("确定", func() {})
	// content2 := container.New(layout.NewBorderLayout(topPos, nil, leftPos, rightBtn),
	// 	topPos, leftPos, rightBtn, fileList)
	// fileListWindow.SetContent(content2)
	// fileListWindow.Resize(fyne.NewSize(600, 400))
	// fileListWindow.CenterOnScreen()
	//content := container.NewVBox(form, infoLabel)

	fileList.OnSelected = func(id widget.ListItemID) {
		input.SetText(listData[id])
	}

	content := container.New(layout.NewBorderLayout(form, infoLabel, nil, searchContent), form, infoLabel, searchContent, fileList)
	myWindow.Resize(fyne.NewSize(600, 500))
	myWindow.SetContent(content)
	myWindow.CenterOnScreen()

	//fileListWindow.Show()
	//fileListWindow.
	myWindow.SetMaster()
	myWindow.ShowAndRun()
	os.Unsetenv("FYNE_FONT")
}
