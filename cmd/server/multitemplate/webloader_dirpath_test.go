package multitemplate

import (
	"testing"
	"testing/fstest"
)

func TestWebloader_pathGeneration(t *testing.T) {
	w := CreateDefaultWebLoader(fstest.MapFS{})

	tests_data := []struct {
		test_name string
		got       string
		want      string
	}{
		{
			test_name: "Template dir path",
			want:      "domain/test_domain/templates",
			got:       w.templatesDirPath("test_domain"),
		}, {
			test_name: "Resources dir path",
			want:      "domain/test_domain/resources",
			got:       w.resourcesDirPath("test_domain"),
		}, {
			test_name: "Pages dir path",
			want:      "domain/test_domain/templates/pages",
			got:       w.pagesDirPath("test_domain"),
		}, {
			test_name: "Layout path",
			want:      "domain/test_domain/templates/layout.html",
			got:       w.defaultLayoutPath("test_domain"),
		},
	}

	for _, td := range tests_data {
		t.Run(td.test_name, func(t *testing.T) {
			assertEqual(t, "path", td.want, td.got)
		})
	}
}
