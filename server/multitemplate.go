package server

import (
	"log"

	"github.com/gin-contrib/multitemplate"
	"github.com/yogusita/to-adhdo/server/funcmap"
)

var pages_base_files = []string{
	"domain/shared/partials/base_styles.html",
	"domain/shared/partials/meta_tags.html",
	"domain/shared/partials/nav.html",
}

type TemplateRenderData struct {
	name   string
	path   string
	layout string
}

var pages_render_data = []TemplateRenderData{
	// wishlist templates
	{name: "wishlist", path: "domain/wishlist/pages/wishlist.html", layout: "domain/wishlist/layouts/default.html"},
	{name: "wishlist-admin", path: "domain/wishlist/pages/wishitem/list.html", layout: "domain/wishlist/layouts/default.html"},
	{name: "wishlist-view", path: "domain/wishlist/pages/wishitem/view.html", layout: "domain/wishlist/layouts/default.html"},
	{name: "wishlist-form", path: "domain/wishlist/pages/wishitem/form.html", layout: "domain/wishlist/layouts/default.html"},

	// articles templates
	{name: "articles-catalog", path: "domain/articles/pages/catalog.html", layout: "domain/articles/layouts/default.html"},
	{name: "articles-view", path: "domain/articles/pages/view.html", layout: "domain/articles/layouts/default.html"},
	{name: "articles-form", path: "domain/articles/pages/form.html", layout: "domain/articles/layouts/default.html"},
	{name: "articles-list", path: "domain/articles/pages/list.html", layout: "domain/articles/layouts/default.html"},

	// tags tempates
	{name: "tags-list", path: "domain/tags/pages/list.html", layout: "domain/tags/layouts/default.html"},
}

var components_render_data = []TemplateRenderData{
	{name: "catalog-list", path: "domain/articles/components/catalog_list.html"},
}

func loadTemplates() multitemplate.Renderer {
	r := multitemplate.NewRenderer()

	fm := funcmap.GetFuncMap()

	component_files := []string{}

	for _, component_data := range components_render_data {
		component_files = append(component_files, component_data.path)
	}

	for _, page_data := range pages_render_data {
		log.Printf("[TEMPLATE][PAGE] Loading template `%s`\n", page_data.name)

		files := []string{
			page_data.path,
		}

		if len(page_data.layout) > 0 {
			files = append([]string{page_data.layout}, files...)
		}

		files = append(files, append(component_files, pages_base_files...)...)

		log.Printf("FILES:\n")
		for i, file := range files {
			log.Printf("[TEMPLATE][PAGE][FILE %d] %s\n", i, file)
		}

		r.AddFromFilesFuncs(
			page_data.name,
			fm,
			files...,
		)

		log.Println("")
	}

	for _, td := range components_render_data {
		log.Printf("[TEMPLATE][COMPONENT] Loading template `%s` from path: %s\n", td.name, td.path)
		r.AddFromFilesFuncs(td.name, fm, td.path)
	}

	return r
}
