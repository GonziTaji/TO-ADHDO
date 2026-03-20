package server

import (
	"fmt"
	"go/format"
	"html/template"
	"io/fs"
	"log"
	"os"
	"path"
	"path/filepath"
	"slices"
	"strings"

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
	{name: "articles-catalog", path: "domain/articles/www/pages/catalog/catalog.html", layout: "domain/articles/www/layouts/default.html"},
	{name: "articles-view", path: "domain/articles/www/pages/view.html", layout: "domain/articles/www/layouts/default.html"},
	{name: "articles-form", path: "domain/articles/www/pages/form.html", layout: "domain/articles/www/layouts/default.html"},
	{name: "articles-list", path: "domain/articles/www/pages/list.html", layout: "domain/articles/www/layouts/default.html"},

	// tags tempates
	{name: "tags-list", path: "domain/tags/pages/list.html", layout: "domain/tags/layouts/default.html"},
}

var components_render_data = []TemplateRenderData{
	{name: "catalog-list", path: "domain/articles/www/components/catalog_list.html"},
}

var resources_extensions = []string{
	".css",
	".js",
	".svg",
	".webp",
}

var templates_extensions = []string{
	".html",
}

var domain_dir = "domain"
var domain_templates_dir = "www"
var domain_pages_dir = "pages"

var domain_page_file = "page.html"
var domain_layout_file = "layout.html"

// Returns whether any of the extensions is equal to the extension of the file_name.
//
// Extensions should be formatted just like the filepath.Ext(s) returns it. I.e.: ".png" instead of "png"
func isAnyExt(file_name string, extensions []string) bool {
	return slices.ContainsFunc(extensions, func(allowed_ext string) bool {
		return strings.EqualFold(allowed_ext, filepath.Ext(file_name))
	})
}

func isTemplate(path string) bool {
	return isAnyExt(path, resources_extensions)
}

func isResource(path string) bool {
	return isAnyExt(path, templates_extensions)
}

type TemplateData struct {
	TemplateName string
	FilePath     string
}

type PageTemplateData struct {
	TemplateName string
	FilePath     string
	LayoutPath string
}

func test_tmpl() {
	tests_subjects := []string{
		"domain/articles/pages/catalog/page.html",
		"domain/shared/partials/nav.html",
		"domain/shared/partials/card/index.html",
		"domain/shared/partials/card/card.js",
	}

	for i, path := range tests_subjects {
		log.Printf("[%d] I: %s\n", i, path)

		out := GetTemplateName(path)

		log.Printf("[%d] O: %s\n\n", i, out)
	}
}

func GetTemplateName(template_file_path string) string {
	rel_path, err := filepath.Rel(domain_dir, template_file_path)

	if err != nil {
		log.Panicf("%s\n", err.Error())
	}

	file_dir, file_name := filepath.Split(rel_path)
	file_dir = filepath.Clean(file_dir)

	if slices.Contains(tmpls_entry_points, file_name) {
		return file_dir
	}

	no_ext_filename := strings.TrimSuffix(file_name, filepath.Ext(file_name))

	return filepath.Join(file_dir, no_ext_filename)
}

func walkDomain(r *multitemplate.Renderer, domain_path string, base_templates_paths []string) error {
	var pages_template_data []PageTemplateData
	var component_paths []TemplateData
	var static_paths []string

	default_layout_path := filepath.Join(domain_path, domain_layout_file)

	// what if no default layout?
	if _, err := os.Stat(default_layout_path); os.IsNotExist(err) {
		log.Printf("WARNING: Domain %s without default layout\n", filepath.Base(domain_path))
	}

	// 1. Get pages with their layout
	pages_dir_path := filepath.Join(domain_path, domain_pages_dir)

	if _, err := os.Stat(default_layout_path); os.IsNotExist(err) {
		// TODO: return error?
		log.Printf("WARNING: Domain %s without pages folder\n", filepath.Base(domain_path))

	}

	pages_dir_entries, err := os.ReadDir(pages_dir_path)

	for _, page_entry := range pages_dir_entries {
		page_dir := filepath.Join(pages_dir_path, page_entry.Name())

		page_data := PageTemplateData{
			TemplateName: page_dir,
			FilePath: filepath.Join(page_dir, domain_page_file),
			LayoutPath: filepath.Join(page_dir, domain_layout_file),
		}

		if _, err := os.Stat(page_data.FilePath); os.IsNotExist(err) {
			log.Printf("WARNING: Page dir %s without page file. Skipping dir\n", page_dir)
			continue
		}

		if _, err := os.Stat(page_data.LayoutPath); os.IsNotExist(err) {
			// use default layout
			page_data.LayoutPath = default_layout_path
		}

		if len(page_data.LayoutPath) == 0 {
			// @review: have a basic fallback? or panic if no fallback?
			log.Printf("WARNING: Page %s without layout in layout-less domain. Skipping page\n", page_dir)
			continue
		}

		pages_template_data = append(pages_template_data, page_data)
	}

	// 2. Get components/partials and static files
	err = filepath.WalkDir(domain_path, func(path string, d fs.DirEntry, err error) error {
		is_template := slices.ContainsFunc(templates_extensions, func(allowed_tmpl_ext string) bool {
			return strings.EqualFold(allowed_tmpl_ext, filepath.Ext(path))
		})

		if is_template {
			domain_pages_dir := filepath.Join(domain_path, domain_pages_dir)

			if strings.Contains(path, domain_pages_dir) {
				// in domain/pages/

				if strings.EqualFold(d.Name(), domain_page_file) || strings.EqualFold(d.Name(), domain_layout_file) {
					// skip pages and layouts
					return nil
				}
			}

			// handle template
			full_template_name := filepath.Join(path, d.Name())

			template_name := filepath.Rel(domain_dir, 

			template_name_parts := filepath.Split()
			filepath.Rel

			if 

			component_paths = append(component_paths, TemplateData{
				TemplateName: 
				FilePath: ,
			})

			return nil
		}

		is_resource := slices.ContainsFunc(resources_extensions, func(allowed_ext string) bool {
			return strings.EqualFold(allowed_ext, filepath.Ext(path))
		})

		if is_resource {
			// TODO: handle resource
		}

		return nil
	})

	if err != nil {
		return err
	}

	// 3. Register pages and components
	// use route as the template, stripping the "domain" part, and the filename name.

	// I.e. domain/articles/pages/components/card/card.html => articles/pages/components/card

	// r.AddFromFilesFuncs(
	// 	page_data.name,
	// 	fm,
	// 	files...,
	// )

	return nil

}

func loadTemplatesAuto() (multitemplate.Renderer, error) {
	r := multitemplate.NewRenderer()

	domains, err := os.ReadDir(domain_dir)

	if err != nil {
		return nil, err
	}

	for _, domain_entry := range domains {
		if !domain_entry.IsDir() {
			continue
		}

		default_layout_path := filepath.Join(domain_dir, domain_entry.Name(), domain_layout_file)

		pages_dir := filepath.Join(domain_dir, domain_entry.Name(), domain_pages_dir)

		pages_entries, err := os.ReadDir(pages_dir)

		if err != nil {
			return nil, err
		}

		for _, page_dir_entry := range pages_entries {
			page_dir_path := filepath.Join(pages_dir, page_dir_entry.Name())

			page_file_path := filepath.Join(page_dir_path, domain_page_file)

			if _, err := os.Stat(); os.IsNotExist(err) {
				log.Print("WARNING: page file not found in %s\n", page_dir_path)
				continue
			}

			page_layout_path := filepath.Join(pages_dir, page_dir_entry.Name(), domain_layout_file)

			if _, err := os.Stat(); os.IsNotExist(err) {
				page_layout_path = default_layout_path
			}

		}

		/*
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

		*/

		err := filepath.WalkDir(pages_dir, func(path string, d fs.DirEntry, err error) error {
			if d.IsDir() {
				return nil
			}

			// is_allowed := slices.ContainsFunc(allowed_extensions, func(allowed_ext string) bool {
			// 	return strings.EqualFold(allowed_ext, filepath.Ext(path))
			// })
			//
			// if !is_allowed {
			// 	return nil
			// }

			return nil
		})

		if err != nil {
			return nil, err
		}
	}

	return r, nil
}

func loadTemplates() multitemplate.Renderer {
	r := multitemplate.NewRenderer()

	fm := funcmap.CreateFuncMap()

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
