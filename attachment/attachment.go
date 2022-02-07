package attachment

import (
	"bytes"
	"fmt"
	"net/http"
	"net/url"
	
	"github.com/gin-gonic/gin"
	"github.com/gladmo/toolset/makelevel"
	"github.com/xuri/excelize/v2"
)

// GinExportExcel gin 附件导出中间件，根据 data 生成表格
func GinExportExcel(c *gin.Context, fileName string, data interface{}) {
	result, err := GenerateExcel(fileName, data)
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		c.Abort()
		return
	}
	
	c.Writer.Header().Set("Content-Disposition", fmt.Sprintf(`attachment; filename="%s.xlsx"`, url.QueryEscape(fileName)))
	c.Data(http.StatusOK, "application/xlsx", result.Bytes())
	c.Abort()
}

// GenerateTable 根据用户自定义的 slice 生成表格类型数据
func GenerateTable(slice interface{}) (data [][]string) {
	ml := makelevel.NewMakeLevel("header")
	lf := ml.MakeLevelStruct(slice)
	
	return lf.BackFillTable(slice)
}

// GenerateExcel 根据data数据结构生成Excel表格
// data 是一个用户自定义结构的 slice
// 需要导出的字段拥有 header tag
//
// type excelData struct {
// 	Username string `header:"username"`
// 	Sex      string `header:"sex"`
// 	Age      string `header:"age"`
// }
func GenerateExcel(sheetName string, data interface{}) (buf bytes.Buffer, err error) {
	table := GenerateTable(data)
	
	f := excelize.NewFile()
	f.SetSheetName("Sheet1", sheetName)
	
	for y, item := range table {
		err = f.SetSheetRow(sheetName, fmt.Sprintf("A%d", y+1), &item)
		if err != nil {
			return
		}
	}
	
	_, err = f.WriteTo(&buf)
	if err != nil {
		return
	}
	
	return
}
