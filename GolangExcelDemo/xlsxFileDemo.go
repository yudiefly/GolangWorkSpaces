package main

import (
	"github.com/tealeg/xlsx"
)

func main() {
	file := xlsx.NewFile()
	sheet, err := file.AddSheet("人员信息收集")
	if err != nil {
		panic(err.Error())
	}
	row := sheet.AddRow()
	cell := row.AddCell()
	cell.Value = "姓名"
	cell = row.AddCell()
	cell.Value = "性别"

	row = sheet.AddRow()
	cell = row.AddCell()
	cell.Value = "张三"
	cell = row.AddCell()
	cell.Value = "男"

	row = sheet.AddRow()
	cell = row.AddCell()
	cell.Value = "李四"
	cell = row.AddCell()
	cell.Value = "女"

	row = sheet.AddRow()
	cell = row.AddCell()
	cell.Value = "王五"
	cell = row.AddCell()
	cell.Value = "男"

	err = file.Save("demo.xlsx")
	if err != nil {
		panic(err.Error())
	}

}
