package articles

import (
	"html/template"
	"io"

	"github.com/yogusita/to-adhdo/domain/articles/model"
)

type Views struct{}

func (v *Views) AsListItem(w io.Writer, article model.Article) error {
	tmpl, err := template.ParseFiles("public/lib/components/articles_list/template.html")

	if err != nil {
		return err
	}

	tmpl.ExecuteTemplate(w, "article_as_list_item", article)

	return nil
}
