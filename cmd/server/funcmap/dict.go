package funcmap

// Usage:
//
// To pass the dictionary { form: .formValues, tags: .tags } to a template block, you do:
//
// {{ template "myBlock" (dict "form" .formValues "tags" .tags) }}
func dict(values ...any) map[string]any {
	if len(values)%2 != 0 {
		panic("dict requires even number of arguments")
	}
	m := make(map[string]any, len(values)/2)
	for i := 0; i < len(values); i += 2 {
		key, ok := values[i].(string)
		if !ok {
			panic("dict keys must be strings")
		}
		m[key] = values[i+1]
	}
	return m
}
