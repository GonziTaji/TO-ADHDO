package multitemplate

import "testing"

func TestWebLoader_getTemplateName(t *testing.T) {
	test_data := []struct {
		test_name     string
		template_path string
		want_name     string
	}{
		{
			test_name:     "Page",
			template_path: "domain/articles/templates/pages/catalog.html",
			want_name:     "articles/pages/catalog",
		}, {
			test_name:     "Dir component",
			template_path: "domain/wishlist/templates/components/card.html",
			want_name:     "wishlist/components/card",
		},
	}

	for _, td := range test_data {
		t.Run(td.test_name, func(t *testing.T) {
			mustHaveName(t, td.template_path, td.want_name)
		})
	}
}
