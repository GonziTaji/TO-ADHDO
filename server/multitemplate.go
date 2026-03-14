package server

import (
	"github.com/gin-contrib/multitemplate"
	"github.com/yogusita/to-adhdo/server/funcmap"
)

func loadTemplates() multitemplate.Renderer {
	r := multitemplate.NewRenderer()

	fm := funcmap.GetFuncMap()

	// Wishlist templates
	r.AddFromFilesFuncs(
		"wishlist",
		fm,
		"domain/wishlist/layouts/default.html",
		"domain/wishlist/pages/wishlist.html",
		"domain/shared/partials/base_styles.html",
		"domain/shared/partials/meta_tags.html",
		"domain/shared/partials/nav.html",
	)

	r.AddFromFilesFuncs(
		"wishlist-admin",
		fm,
		"domain/wishlist/layouts/default.html",
		"domain/wishlist/pages/wishitem/list.html",
		"domain/shared/partials/base_styles.html",
		"domain/shared/partials/meta_tags.html",
		"domain/shared/partials/nav.html",
	)

	r.AddFromFilesFuncs(
		"wishlist-view",
		fm,
		"domain/wishlist/layouts/default.html",
		"domain/wishlist/pages/wishitem/view.html",
		"domain/shared/partials/base_styles.html",
		"domain/shared/partials/meta_tags.html",
		"domain/shared/partials/nav.html",
	)

	r.AddFromFilesFuncs(
		"wishlist-form",
		fm,
		"domain/wishlist/layouts/default.html",
		"domain/wishlist/pages/wishitem/form.html",
		"domain/shared/partials/base_styles.html",
		"domain/shared/partials/meta_tags.html",
		"domain/shared/partials/nav.html",
	)

	// Articles templates
	r.AddFromFilesFuncs(
		"articles-catalog",
		fm,
		"domain/articles/layouts/default.html",
		"domain/articles/pages/catalog.html",
		"domain/shared/partials/base_styles.html",
		"domain/shared/partials/meta_tags.html",
		"domain/shared/partials/nav.html",
	)

	r.AddFromFilesFuncs(
		"articles-view",
		fm,
		"domain/articles/layouts/default.html",
		"domain/articles/pages/view.html",
		"domain/shared/partials/base_styles.html",
		"domain/shared/partials/meta_tags.html",
		"domain/shared/partials/nav.html",
	)

	r.AddFromFilesFuncs(
		"articles-form",
		fm,
		"domain/articles/layouts/default.html",
		"domain/articles/pages/form.html",
		"domain/shared/partials/base_styles.html",
		"domain/shared/partials/meta_tags.html",
		"domain/shared/partials/nav.html",
	)

	r.AddFromFilesFuncs(
		"articles-list",
		fm,
		"domain/articles/layouts/default.html",
		"domain/articles/pages/list.html",
		"domain/articles/static/templates/articles_list/template.html",
		"domain/shared/partials/base_styles.html",
		"domain/shared/partials/meta_tags.html",
		"domain/shared/partials/nav.html",
	)

	// Tags templates
	r.AddFromFilesFuncs(
		"tags-list",
		fm,
		"domain/tags/layouts/default.html",
		"domain/tags/pages/list.html",
		"domain/shared/partials/base_styles.html",
		"domain/shared/partials/meta_tags.html",
		"domain/shared/partials/nav.html",
	)

	return r
}
