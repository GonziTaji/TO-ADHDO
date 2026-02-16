package funcmap

import "html/template"

func GetFuncMap() template.FuncMap {
	return template.FuncMap{
		"dict":       dict,
		"add":        add,
		"format_clp": format_clp,
	}
}
