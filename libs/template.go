package libs

import (
	"html/template"
	"os"
)

// 解析模板，生成临时命令文件
func ParseTemplate(tmplFilename, tempFilename string, data interface{}) error {
	tmpl, err := template.ParseFiles(tmplFilename)
	if err != nil {
		return err
	}

	tempFile, err := os.Create(tempFilename)
	if err != nil {
		return err
	}
	defer tempFile.Close()

	err = tmpl.Execute(tempFile, data)
	if err != nil {
		return err
	}

	return nil
}
