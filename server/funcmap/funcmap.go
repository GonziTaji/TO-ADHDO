package funcmap

import "html/template"

type TemplateFuncMaps struct {
	funcMap template.FuncMap
}

func CreateFuncMap() *TemplateFuncMaps {
	return &TemplateFuncMaps{
		funcMap: template.FuncMap{
			"dict":       dict,
			"add":        add,
			"format_clp": format_clp,
			"first":      first,
		},
	}
}

// It returns a funcmap with the `resource` function pointing to the proper domain URI
func (t *TemplateFuncMaps) FuncMap(domain string) template.FuncMap {
	fm := t.funcMap
	fm["resource"] = 

	return t.funcMap
}
