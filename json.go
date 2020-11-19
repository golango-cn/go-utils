package go_utils

import (
	"encoding/json"
	"strings"
)

// 移除Json字段
func RemoveJsonField(fields []string, mp map[string]interface{}) map[string]interface{} {

	if mp == nil {
		return mp
	}
	if len(fields) == 0 {
		return mp
	}

	for k, v := range mp {
		inarray := false
		for _, field := range fields {
			if strings.EqualFold(field, k) {
				inarray = true
				break
			}
		}
		if inarray {
			delete(mp, k)
		}
		switch x := v.(type) {
		case map[string]interface{}:
			x = RemoveJsonField(fields, x)
		}
	}

	return mp

}

func RemoveJsonStringField(fields []string, str string) string {

	if len(fields) == 0 {
		return str
	}
	if len(str) == 0 {
		return str
	}

	var mp map[string]interface{}
	err := json.Unmarshal([]byte(str), &mp)
	if err != nil {
		return str
	}

	mp = RemoveJsonField(fields, mp)
	b, _ := json.Marshal(mp)

	return string(b)
}
