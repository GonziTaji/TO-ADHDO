package articles

import (
	"html/template"
	"io"
)

type Views struct{}

func (v *Views) AsListItem(w io.Writer, article Article) error {
	tmpl, err := template.ParseFiles("public/lib/components/articles_list/template.html")

	if err != nil {
		return err
	}

	tmpl.ExecuteTemplate(w, "article_as_list_item", article)

	return nil
}
