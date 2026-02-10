package funcmap

import "html/template"

func GetFuncMap() template.FuncMap {
	return template.FuncMap{
		"dict": dict,
		"add":  add,
	}
}
