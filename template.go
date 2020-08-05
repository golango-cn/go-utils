package go_utils

import (
	"bytes"
	"html/template"
)

// 转换模板
func ParseTemplate(data interface{}, templateFile string) (string, error) {

	t, err := template.ParseFiles(templateFile)
	if err != nil {
		return "", err
	}

	buf := new(bytes.Buffer)
	if err = t.Execute(buf, data); err != nil {
		return "", err
	}
	return buf.String(), nil

}
