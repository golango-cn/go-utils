package go_utils

import (
	"bytes"
	"errors"
	"html/template"
)

// 转换模板
func ParseTemplate(name string, funcMap template.FuncMap, data interface{}, templateFile string) (string, error) {

	if len(name) == 0 {
		return "", errors.New("Name cannot be empty")
	}
	t := template.New(name)
	if funcMap != nil {
		t = t.Funcs(funcMap)
	}
	t, err := t.ParseFiles(templateFile)
	if err != nil {
		return "", err
	}
	buf := new(bytes.Buffer)
	if err = t.Execute(buf, data); err != nil {
		return "", err
	}
	return buf.String(), nil

}
