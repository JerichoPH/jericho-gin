package controllers

import (
	"fmt"
	"os"

	"jericho-go/tools"
	"jericho-go/wrongs"

	"github.com/gin-gonic/gin"
)

// CommandController 控制台控制器
type CommandController struct{}

// NewCommandController 构造函数
func NewCommandController() *CommandController {
	return &CommandController{}
}

// ExcelHelperDemo 列表
func (receiver CommandController) ExcelHelperDemo(ctx *gin.Context) {
	dir := os.Getenv("PWD")
	operation := ctx.Query("operation")
	excelName := ctx.Query("excel_name")

	if operation == "read" {
		excelReader := (&tools.ExcelReader{}).
			OpenFile(fmt.Sprintf("%s/static/%s.xlsx", dir, excelName)).
			SetSheetName("Sheet1").
			ReadTitle().
			Read()

		fmt.Println(excelReader.GetTitle(), excelReader.ToList(), excelReader.ToMap())
		fmt.Println("----------")
		fmt.Println(excelReader.ToDataFrameDefaultType().Records())

		ctx.JSON(
			tools.NewCorrectWithGinContext("", ctx).
				Datum(
					map[string]any{
						"title":       excelReader.GetTitle(),
						"listByQuery": excelReader.ToList(),
						"map":         excelReader.ToMap(),
						"dataframe":   excelReader.ToDataFrameDefaultType().Maps(),
					},
				).ToGinResponse(),
		)
	} else if operation == "write" {
		// 写入Excel
		// 设置表头
		titleRow := new(tools.ExcelRow).
			SetRowNumber(1).
			SetCells([]*tools.ExcelCell{
				new(tools.ExcelCell).SetContent("姓名").SetFontColor("#FF0000", true),
				new(tools.ExcelCell).SetContent("年龄"),
				new(tools.ExcelCell).SetContent("性别"),
			})
		err := (&tools.ExcelWriter{}).
			Init(fmt.Sprintf("%s/static/%s.xlsx", dir, excelName)).
			ActiveSheetByIndex(0).
			SetRows([]*tools.ExcelRow{titleRow}).
			Save()
		if err != nil {
			wrongs.ThrowForbidden(fmt.Sprintf("保存文件失败：%s", err.Error()))
		}
	}
}

// InitData 初始化数据
func (receiver CommandController) InitData(ctx *gin.Context) {
}
