package multitemplate

import (
	"fmt"
	"slices"
	"sort"
	"strings"
	"testing"
	"testing/fstest"
)

// makeLoader builds a WebLoader backed by a MapFS with the default config.
func makeLoader(fsys fstest.MapFS) *WebLoader {
	return CreateWebLoader(fsys, DefaultWebConfig())
}

// makeLoaderWithCfg builds a WebLoader backed by a MapFS with a custom config.
func makeLoaderWithCfg(fsys fstest.MapFS, c WebConfig) *WebLoader {
	return CreateWebLoader(fsys, c)
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

func sortResources(res []string) {
	sort.Strings(res)
}

// ---------------------------------------------------------------------------
// Test 1 – Components vs Pages separation
// ---------------------------------------------------------------------------

func TestWebLoader_getDomainData_ComponentsVsPagesSeparation(t *testing.T) {
	fsys := fstest.MapFS{
		"domain/app/pages/home/page.html":   file("page"),
		"domain/app/pages/home/layout.html": file("layout"),
		"domain/app/partials/header.html":   file("header"),
	}

	w := makeLoader(fsys)
	data, err := w.getDomainData("app")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	sortPages(data.Pages)

	// NOTE: per-page layout.html inside pages/ uses a different map key than
	// page.html (key = "app/pages/home/layout" vs "app/pages/home"), so they
	// don't merge into the same PageTemplateData. This results in two page
	// entries: one real page entry and one "phantom" entry for the layout file
	// with an empty TemplateName (only LayoutPath is set, TemplateName is never
	// assigned for the layout branch).
	if len(data.Pages) != 1 {
		t.Fatalf("expected 1 page entry, got %d: %v", len(data.Pages), data.Pages)
	}

	// The real page entry has a TemplateName set.
	var realPage *PageTemplateData
	for i := range data.Pages {
		if data.Pages[i].TemplateName == "app/pages/home" {
			realPage = &data.Pages[i]
			break
		}
	}
	if realPage == nil {
		t.Fatalf("expected page with TemplateName %q, pages: %v", "app/pages/home", data.Pages)
	}
	if realPage.FilePath != "domain/app/pages/home/page.html" {
		t.Errorf("page FilePath: want %q, got %q", "domain/app/pages/home/page.html", realPage.FilePath)
	}

	// page.html and layout.html inside pages/ must NOT appear as components.
	sortComponents(data.Components)
	for _, comp := range data.Components {
		if comp.FilePath == "domain/app/pages/home/page.html" {
			t.Errorf("page.html should NOT appear as a component")
		}
		if comp.FilePath == "domain/app/pages/home/layout.html" {
			t.Errorf("layout.html should NOT appear as a component")
		}
	}

	// partials/header.html must appear as a component.
	found := slices.ContainsFunc(data.Components, func(c TemplateData) bool {
		return c.FilePath == "domain/app/partials/header.html"
	})
	if !found {
		t.Errorf("expected domain/app/partials/header.html in components, got %v", data.Components)
	}
}

// ---------------------------------------------------------------------------
// Test 2 – Default layout fallback
// ---------------------------------------------------------------------------

func TestWebLoader_getDomainData_DefaultLayoutFallback(t *testing.T) {
	fsys := fstest.MapFS{
		// Default layout at domain root.
		"domain/app/pages/layout.html":    file("default layout"),
		"domain/app/pages/home/page.html": file("home page"),
		// No per-page layout.html for home.
	}

	w := makeLoader(fsys)
	data, err := w.getDomainData("app")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	sortPages(data.Pages)

	var homePage *PageTemplateData
	for i, p := range data.Pages {
		if p.TemplateName == "app/pages/home" {
			homePage = &data.Pages[i]
			break
		}
	}
	if homePage == nil {
		t.Fatalf("expected page app/pages/home, pages: %v", data.Pages)
	}

	wantLayout := "domain/app/pages/layout.html"
	if homePage.LayoutPath != wantLayout {
		t.Errorf("LayoutPath: want %q, got %q", wantLayout, homePage.LayoutPath)
	}
}

// ---------------------------------------------------------------------------
// Test 3 – Resources collected; templates not misclassified
// ---------------------------------------------------------------------------

func TestWebLoader_getDomainData_ResourcesCollected(t *testing.T) {
	fsys := fstest.MapFS{
		"domain/app/assets/site.css":      file("/* css */"),
		"domain/app/assets/app.js":        file("// js"),
		"domain/app/partials/header.html": file("<header>"),
	}

	w := makeLoader(fsys)
	data, err := w.getDomainData("app")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	sortResources(data.Resources)

	wantResources := []string{
		"domain/app/assets/app.js",
		"domain/app/assets/site.css",
	}
	if len(data.Resources) != len(wantResources) {
		t.Fatalf("resources: want %v, got %v", wantResources, data.Resources)
	}
	for i, want := range wantResources {
		if data.Resources[i] != want {
			t.Errorf("resource[%d]: want %q, got %q", i, want, data.Resources[i])
		}
	}

	// header.html must appear as a component, not a resource
	sortComponents(data.Components)
	if len(data.Components) != 1 {
		t.Fatalf("components: want 1, got %d: %v", len(data.Components), data.Components)
	}
	if data.Components[0].FilePath != "domain/app/partials/header.html" {
		t.Errorf("component FilePath: want %q, got %q",
			"domain/app/partials/header.html", data.Components[0].FilePath)
	}

	// resources must not include templates
	for _, res := range data.Resources {
		if res == "domain/app/partials/header.html" {
			t.Error("header.html must not appear as a resource")
		}
	}
}

// ---------------------------------------------------------------------------
// Test 4 – Case-insensitive extension handling
// ---------------------------------------------------------------------------

func TestWebLoader_getDomainData_CaseInsensitiveExtensions(t *testing.T) {
	fsys := fstest.MapFS{
		"domain/app/partials/HEADER.HTML": file("<header>"),
		"domain/app/assets/SITE.CSS":      file("/* css */"),
	}

	w := makeLoader(fsys)
	data, err := w.getDomainData("app")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	foundTemplate := slices.ContainsFunc(data.Components, func(c TemplateData) bool {
		return c.FilePath == "domain/app/partials/HEADER.HTML"
	})
	if !foundTemplate {
		t.Errorf("expected HEADER.HTML to be detected as a template (case-insensitive), components: %v", data.Components)
	}

	foundResource := slices.Contains(data.Resources, "domain/app/assets/SITE.CSS")
	if !foundResource {
		t.Errorf("expected SITE.CSS to be detected as a resource (case-insensitive), resources: %v", data.Resources)
	}
}

// ---------------------------------------------------------------------------
// Test 5 – Injected config is respected
// ---------------------------------------------------------------------------

func TestWebLoader_getDomainData_ConfigInjection(t *testing.T) {
	// Use a non-default DomainDirName and PagesDirName / filenames.
	c := WebConfig{
		DomainDirName:       "tenants",
		PagesDirName:        "views",
		PageFileName:        "index.html",
		LayoutFileName:      "base.html",
		TemplatesExtensions: []string{".html"},
		ResourcesExtensions: []string{".css"},
	}

	fsys := fstest.MapFS{
		"tenants/shop/views/catalog/index.html": file("catalog page"),
		"tenants/shop/views/catalog/base.html":  file("catalog layout"),
		"tenants/shop/partials/nav.html":        file("nav"),
		"tenants/shop/assets/main.css":          file("/* css */"),
	}

	w := makeLoaderWithCfg(fsys, c)
	data, err := w.getDomainData("shop")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	sortPages(data.Pages)

	if len(data.Pages) != 1 {
		t.Fatalf("expected 2 page entries, got %d: %v", len(data.Pages), data.Pages)
	}

	var catalogPage *PageTemplateData
	for i := range data.Pages {
		if data.Pages[i].TemplateName == "shop/views/catalog" {
			catalogPage = &data.Pages[i]
			break
		}
	}
	if catalogPage == nil {
		t.Fatalf("expected page shop/views/catalog, pages: %v", data.Pages)
	}
	if catalogPage.FilePath != "tenants/shop/views/catalog/index.html" {
		t.Errorf("page FilePath: want %q, got %q",
			"tenants/shop/views/catalog/index.html", catalogPage.FilePath)
	}

	// nav.html should be a component
	foundNav := slices.ContainsFunc(data.Components, func(c TemplateData) bool {
		return c.FilePath == "tenants/shop/partials/nav.html"
	})
	if !foundNav {
		t.Errorf("expected nav.html as component, components: %v", data.Components)
	}

	// main.css should be a resource
	if !slices.Contains(data.Resources, "tenants/shop/assets/main.css") {
		t.Errorf("expected main.css as resource, resources: %v", data.Resources)
	}
}

// ---------------------------------------------------------------------------
// Test 6 – Golden test: full DomainTemplatesData for a representative tree
// ---------------------------------------------------------------------------

// TestWebLoader_getDomainData_Golden builds a realistic domain tree covering
// pages (with and without per-page layout), components/partials, and static
// resources, then asserts the full DomainTemplatesData exactly.
//
// NOTE: per-page layout detection uses the same key derived from
// getTemplateNameWithBase. Because "layout.html" is not in tmpls_entry_points,
// its key differs from the page key (e.g. "app/pages/home/layout" vs
// "app/pages/home"). As a result, per-page layouts produce a separate pages_by_name
// entry (FilePath empty, LayoutPath set) which then also appears as a Page entry
// with the default layout. This is the current documented behavior; a future
// refactor can align keys explicitly.
func TestWebLoader_getDomainData_Golden(t *testing.T) {
	fsys := fstest.MapFS{
		// Default layout at domain root
		"domain/app/pages/layout.html": file("default layout"),

		// Page 1: home — no per-page layout, falls back to default
		"domain/app/pages/home/page.html": file("home"),

		// Page 2: about — has a per-page layout.html (separate key, see note above)
		"domain/app/pages/about/page.html":   file("about"),
		"domain/app/pages/about/layout.html": file("about layout"),

		// Components/partials
		"domain/app/partials/nav.html":    file("nav"),
		"domain/app/partials/footer.html": file("footer"),

		// Resources
		"domain/app/assets/site.css": file("/* css */"),
		"domain/app/assets/app.js":   file("// js"),
		"domain/app/assets/logo.svg": file("<svg/>"),
	}

	w := makeLoader(fsys)
	data, err := w.getDomainData("app")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// ---- Pages ----
	// Expected page entries:
	//   "app/pages/home"         → page.html, default layout
	//   "app/pages/about"        → page.html, default layout (per-page layout
	//                              lands under a different key, see note above)
	sortPages(data.Pages)

	type wantPage struct {
		TemplateName string
		FilePath     string
		LayoutPath   string
	}
	// Sort order: empty string sorts before any non-empty string.
	// The phantom entry (layout.html → key "app/pages/about/layout") never gets
	// its TemplateName set (only LayoutPath is populated), so TemplateName = "".
	// Since LayoutPath != "", the default-layout fallback is not applied.
	wantPages := []wantPage{
		{
			TemplateName: "app/pages/about",
			FilePath:     "domain/app/pages/about/page.html",
			LayoutPath:   "domain/app/pages/about/layout.html",
		},
		{
			TemplateName: "app/pages/home",
			FilePath:     "domain/app/pages/home/page.html",
			LayoutPath:   "domain/app/pages/layout.html",
		},
	}

	if len(data.Pages) != len(wantPages) {
		t.Fatalf("pages count: want %d, got %d\n  got: %v", len(wantPages), len(data.Pages), data.Pages)
	}
	for i, want := range wantPages {
		got := data.Pages[i]
		if got.TemplateName != want.TemplateName {
			t.Errorf("pages[%d].TemplateName: want %q, got %q", i, want.TemplateName, got.TemplateName)
		}
		if got.FilePath != want.FilePath {
			t.Errorf("pages[%d].FilePath: want %q, got %q", i, want.FilePath, got.FilePath)
		}
		if got.LayoutPath != want.LayoutPath {
			t.Errorf("pages[%d].LayoutPath: want %q, got %q", i, want.LayoutPath, got.LayoutPath)
		}
	}

	// ---- Components ----
	// domain/app/layout.html is NOT under pages/, so it is a component.
	// partials/nav.html and partials/footer.html are components.
	// pages/home/page.html and pages/about/page.html are NOT components.
	// pages/about/layout.html is NOT a component (it's under pages/).
	sortComponents(data.Components)

	wantComponents := []TemplateData{
		{TemplateName: "app/partials/footer", FilePath: "domain/app/partials/footer.html"},
		{TemplateName: "app/partials/nav", FilePath: "domain/app/partials/nav.html"},
	}

	if len(data.Components) != len(wantComponents) {
		t.Fatalf("components count: want %d, got %d\n  got: %v",
			len(wantComponents), len(data.Components), data.Components)
	}
	for i, want := range wantComponents {
		got := data.Components[i]
		if got.TemplateName != want.TemplateName {
			t.Errorf("components[%d].TemplateName: want %q, got %q", i, want.TemplateName, got.TemplateName)
		}
		if got.FilePath != want.FilePath {
			t.Errorf("components[%d].FilePath: want %q, got %q", i, want.FilePath, got.FilePath)
		}
	}

	// ---- Resources ----
	sortResources(data.Resources)

	wantResources := []string{
		"domain/app/assets/app.js",
		"domain/app/assets/logo.svg",
		"domain/app/assets/site.css",
	}

	if len(data.Resources) != len(wantResources) {
		t.Fatalf("resources count: want %d, got %d\n  got: %v",
			len(wantResources), len(data.Resources), data.Resources)
	}
	for i, want := range wantResources {
		if data.Resources[i] != want {
			t.Errorf("resources[%d]: want %q, got %q", i, want, data.Resources[i])
		}
	}

	// ---- Domain name ----
	if data.Name != "app" {
		t.Errorf("Name: want %q, got %q", "app", data.Name)
	}
}

func mustHaveName(input_path string, expected_name string) error {
	w := makeLoader(fstest.MapFS{})

	if out := w.getTemplateName(input_path); !strings.EqualFold(out, expected_name) {
		return fmt.Errorf("Invalid template name. Expected: %s; Got: %s ", expected_name, out)
	}

	return nil
}

func TestWebLoader_getTemplateName_Page(t *testing.T) {
	err := mustHaveName(
		"domain/articles/pages/catalog/page.html",
		"articles/pages/catalog",
	)

	if err != nil {
		t.Error(err.Error())
	}
}

func TestWebLoader_getTemplateName_DirComponent(t *testing.T) {
	err := mustHaveName(
		"domain/wishlist/components/card/index.html",
		"wishlist/components/card",
	)

	if err != nil {
		t.Error(err.Error())
	}
}

func TestWebLoader_getTemplateName_FileComponent(t *testing.T) {
	err := mustHaveName(
		"domain/shared/partials/nav.html",
		"shared/partials/nav",
	)

	if err != nil {
		t.Error(err.Error())
	}
}
