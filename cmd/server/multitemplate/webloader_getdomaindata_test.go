package multitemplate

import (
	"fmt"
	"path"
	"sort"
	"strings"
	"testing"
	"testing/fstest"
)

// makeLoader builds a WebLoader backed by a MapFS with the default config.
func makeLoader(fsys fstest.MapFS) *WebLoader {
	return CreateDefaultWebLoader(fsys)
}

func file(content string) *fstest.MapFile {
	return &fstest.MapFile{Data: []byte(content)}
}

// Sorting helpers — getDomainData builds Pages from a map, so order is not
// guaranteed. Sort all slices before making assertions.

func sortPages(pages []PageTemplateData) {
	sort.Slice(pages, func(i, j int) bool {
		return pages[i].TemplateName < pages[j].TemplateName
	})
}

func sortComponents(comps []TemplateData) {
	sort.Slice(comps, func(i, j int) bool {
		return comps[i].TemplateName < comps[j].TemplateName
	})
}

func assertEqual[T comparable](t *testing.T, name string, want, got T) {
	t.Helper()

	if want != got {
		t.Errorf("%s: expected %v, got %v", name, want, got)
	}
}

func assertPageEqual(t *testing.T, want, got PageTemplateData) {
	t.Helper()

	assertEqual(t, "TemplateName", want.TemplateName, got.TemplateName)
	assertEqual(t, "FilePath", want.FilePath, got.FilePath)
	assertEqual(t, "LayoutPath", want.LayoutPath, got.LayoutPath)
}

func assertComponentEqual(t *testing.T, want, got TemplateData) {
	t.Helper()

	assertEqual(t, "TemplateName", want.TemplateName, got.TemplateName)
	assertEqual(t, "FilePath", want.FilePath, got.FilePath)
}

func validateGetTemplateName(t *testing.T, fsys fstest.MapFS, cfg WebLoaderConfig, domain_name string, want DomainTemplatesData) {
	t.Helper()

	w := CreateDefaultWebLoader(fsys)
	data, err := w.LoadTemplates(domain_name)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	sortPages(data.Pages)
	sortComponents(data.Components)

	sortPages(want.Pages)
	sortComponents(want.Components)

	assertEqual(t, "Domain name", want.DomainName, data.DomainName)
	assertEqual(t, "Pages count", len(want.Pages), len(data.Pages))
	assertEqual(t, "Components count", len(want.Components), len(data.Components))

	for i, want_page := range want.Pages {
		t.Run(fmt.Sprintf("Page %s", want_page.TemplateName), func(t *testing.T) {
			assertPageEqual(t, want_page, data.Pages[i])
		})
	}

	for i, want_component := range want.Components {
		t.Run(fmt.Sprintf("Component %s", want_component.TemplateName), func(t *testing.T) {
			assertComponentEqual(t, want_component, data.Components[i])
		})
	}
}

const baseTemplatePath = "domain/app/www/templates"

func joinBase(p string) string { return path.Join(baseTemplatePath, p) }

func TestWebLoader_getTemplatesFromDomain(t *testing.T) {
	tests := []struct {
		name string
		fsys fstest.MapFS
		want DomainTemplatesData
	}{
		{
			name: "Page with default layout",
			fsys: fstest.MapFS{
				joinBase("pages/layout.html"):     file("default layout"),
				joinBase("pages/home/index.html"): file("home page"),
			},
			want: DomainTemplatesData{
				DomainName: "app",
				Pages: []PageTemplateData{
					{
						LayoutPath:   joinBase("pages/layout.html"),
						FilePath:     joinBase("pages/home/index.html"),
						TemplateName: "www/pages/home",
					},
				},
			},
		},
		{
			name: "Page with default layout",
			fsys: fstest.MapFS{
				joinBase("pages/layout.html"):     file("default layout"),
				joinBase("pages/home/index.html"): file("home page"),
			},
			want: DomainTemplatesData{
				DomainName: "app",
				Pages: []PageTemplateData{
					{
						LayoutPath:   joinBase("pages/layout.html"),
						FilePath:     joinBase("pages/home/index.html"),
						TemplateName: "www/pages/home",
					},
				},
			},
		},
		{
			name: "Page with own layout",
			fsys: fstest.MapFS{
				joinBase("pages/layout.html"):      file("default layout"),
				joinBase("pages/home/layout.html"): file("home layout"),
				joinBase("pages/home/index.html"):  file("home page"),
			},
			want: DomainTemplatesData{
				DomainName: "app",
				Pages: []PageTemplateData{
					{
						LayoutPath:   joinBase("pages/layout.html"),
						FilePath:     joinBase("pages/home/index.html"),
						TemplateName: "www/pages/home",
					},
				},
			},
		},
		{
			name: "Components",
			fsys: fstest.MapFS{
				joinBase("pages/home/components/jumbotron/index.html"): file("a jumbotron"),
				joinBase("components/nav/index.html"):                  file("nav component"),
			},
			want: DomainTemplatesData{
				DomainName: "app",
				Components: []TemplateData{
					{
						FilePath:     joinBase("pages/home/components/jumbotron/index.html"),
						TemplateName: joinBase("www/pages/home/components/jumbotron"),
					},
					{
						FilePath:     joinBase("components/nav/index.html"),
						TemplateName: joinBase("www/components/nav"),
					},
				},
			},
		},
	}

	for _, test_data := range tests {
		t.Run(test_data.name, func(t *testing.T) {
			validateGetTemplateName(t, test_data.fsys, DefaultWebConfig(), "app", test_data.want)
		})
	}
}

func TestWebLoader_getTemplatesFromDomain_CustomConfig(t *testing.T) {
	t.Run("Custom template extension", func(t *testing.T) {
		cfg := DefaultWebConfig()
		cfg.TemplatesExtensions = []string{".gohtml", ".tmpl"}

		test_data := []struct {
			test_name     string
			template_path string
			want_name     string
		}{
			{
				test_name:     "page",
				template_path: "domain/somedomain/pages/home/index.gohtml",
				want_name:     "shared/partials/nav",
			},
			{
				test_name:     "file component",
				template_path: "domain/shared/partials/nav.gohtml",
				want_name:     "shared/partials/nav",
			},
			{
				test_name:     "dir component",
				template_path: "domain/somedomain/pages/home/nav/nav.gohtml",
				want_name:     "shared/partials/nav",
			},
			{
				test_name:     "page with second ext",
				template_path: "domain/somedomain/pages/home/nav/nav.tmpl",
				want_name:     "shared/partials/nav",
			},
		}

		for _, td := range test_data {
			t.Run(td.test_name, func(t *testing.T) {
				mustHaveName(t, td.template_path, td.want_name)
			})
		}

		t.Run("Ignore invalid template ext", func(t *testing.T) {
			fsys := fstest.MapFS{
				"domain/tmpl/pages/home/index.html":                file(""),
				"domain/tmpl/pages/home/components/something.html": file(""),
				"domain/tmpl/components/something_else.html":       file(""),
			}

			validateGetTemplateName(t, fsys, cfg, "tmpl", DomainTemplatesData{
				DomainName: "tmpl",
				Pages:      []PageTemplateData{},
				Components: []TemplateData{},
			})
		})
	})

	// TODO:
	// t.Run("Custom domain dir", func(t *testing.T) {
	// 	cfg := DefaultWebConfig()
	// 	cfg.DomainDirName = "superapp"
	// })
	//
	// one test for each case?
	// cfg.PagesDirName = "views"
	//
	// cfg.LayoutFileName = "base.html" // change extension if needed
	//
	// cfg.TemplatesDirName = "webapp" // change the name or remove this attribute
}

func mustHaveName(t *testing.T, input_path string, expected_name string) {
	t.Helper()

	w := makeLoader(fstest.MapFS{})

	if out := w.getTemplateName(input_path); !strings.EqualFold(out, expected_name) {
		t.Errorf("Invalid template name. Expected: %s; Got: %s ", expected_name, out)
	}
}
