package multitemplate

import (
	"io/fs"
	"log"
	"path/filepath"
	"slices"
	"strings"

	"github.com/gin-contrib/multitemplate"
	"github.com/yogusita/to-adhdo/server/funcmap"
)

type WebConfig struct {
	PageFileName   string
	LayoutFileName string

	DomainDirName    string
	TemplatesDirName string
	PagesDirName     string

	ResourcesExtensions []string
	TemplatesExtensions []string
}

var cfg = WebConfig{
	PageFileName:   "page.html",
	LayoutFileName: "layout.html",

	DomainDirName:    "domain",
	TemplatesDirName: "www",
	PagesDirName:     "pages",

	ResourcesExtensions: []string{
		".css",
		".js",
		".svg",
		".webp",
	},

	TemplatesExtensions: []string{
		".html",
	},
}

func DefaultWebConfig() WebConfig {
	return cfg
}

// Returns whether any of the extensions is equal to the extension of the file_name.
//
// Extensions should be formatted just like the filepath.Ext(s) returns it. I.e.: ".png" instead of "png"
func isAnyExt(file_name string, extensions []string) bool {
	return slices.ContainsFunc(extensions, func(allowed_ext string) bool {
		return strings.EqualFold(allowed_ext, filepath.Ext(file_name))
	})
}

type TemplateData struct {
	TemplateName string
	FilePath     string
}

type PageTemplateData struct {
	TemplateName string
	FilePath     string
	LayoutPath   string
}

type DomainTemplatesData struct {
	Name       string
	Components []TemplateData
	Pages      []PageTemplateData
	Resources  []string
}

type WebLoader struct {
	fs  fs.FS
	cfg WebConfig
}

func CreateWebLoader(fsys fs.FS, cfg WebConfig) *WebLoader {
	return &WebLoader{
		fs:  fsys,
		cfg: cfg,
	}
}

func (w *WebLoader) Load() (multitemplate.Renderer, error) {
	r := multitemplate.NewRenderer()
	fm := funcmap.CreateFuncMap()

	domains, err := fs.ReadDir(w.fs, w.cfg.DomainDirName)

	if err != nil {
		return nil, err
	}

	for _, domain_entry := range domains {
		if !domain_entry.IsDir() {
			continue
		}

		domain_name := domain_entry.Name()

		domain_data, err := w.getDomainData(domain_name)

		if err != nil {
			return nil, err
		}

		log.Printf("Processing domain folder %s\n", domain_data.Name)
		log.Printf("Data: \n\t%v\n", domain_data.Name)

		w.loadDomainFiles(&r, fm, domain_data)
	}

	return r, nil
}

func (w *WebLoader) getDomainData(domain_name string) (DomainTemplatesData, error) {
	data := DomainTemplatesData{
		Name: domain_name,
	}

	pages_by_name := make(map[string]PageTemplateData)

	domain_path := filepath.Join(w.cfg.DomainDirName, domain_name)
	www_path := filepath.Join(domain_path, w.cfg.TemplatesDirName)
	pages_dir_path := filepath.Join(www_path, w.cfg.PagesDirName) + string(filepath.Separator)
	default_layout_path := filepath.Join(pages_dir_path, w.cfg.LayoutFileName)

	err := fs.WalkDir(w.fs, www_path, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if d.IsDir() {
			return nil
		}

		if w.isTemplate(d.Name()) {
			template_name := w.getTemplateName(path)

			if strings.HasPrefix(path, pages_dir_path) {
				if strings.EqualFold(path, default_layout_path) {
					return nil
				}

				if strings.EqualFold(d.Name(), w.cfg.PageFileName) {
					tmp := pages_by_name[template_name]
					tmp.TemplateName = template_name
					tmp.FilePath = path
					pages_by_name[template_name] = tmp

					return nil
				}

				if strings.EqualFold(d.Name(), w.cfg.LayoutFileName) {
					tmp := pages_by_name[template_name]
					tmp.LayoutPath = path
					pages_by_name[template_name] = tmp

					return nil
				}
			}

			data.Components = append(data.Components, TemplateData{
				TemplateName: template_name,
				FilePath:     path,
			})

			return nil
		}

		if w.isResource(d.Name()) {
			data.Resources = append(data.Resources, path)
		}

		return nil
	})

	if err != nil {
		return data, err
	}

	for _, page_data := range pages_by_name {
		if len(page_data.LayoutPath) == 0 {
			page_data.LayoutPath = default_layout_path
		}

		data.Pages = append(data.Pages, page_data)
	}

	return data, nil
}

func (w *WebLoader) loadDomainFiles(r *multitemplate.Renderer, fm *funcmap.TemplateFuncMaps, domain_data DomainTemplatesData) {
	domain_fm := fm.FuncMap(domain_data.Name)

	var components_files []string

	for _, component_data := range domain_data.Components {
		components_files = append(components_files, component_data.FilePath)

		(*r).AddFromFilesFuncs(component_data.TemplateName, domain_fm, component_data.FilePath)
	}

	log.Printf("Component files: %v\n", components_files)

	for _, page_data := range domain_data.Pages {
		log.Printf("Processing page %s\n", page_data.TemplateName)
		log.Printf("Data: \n\t%v\n", page_data)

		files := []string{
			page_data.LayoutPath,
			page_data.FilePath,
		}

		files = append(files, components_files...)

		log.Println("Files:")
		for _, f := range files {
			log.Printf("\t- %s\n", f)
		}

		(*r).AddFromFilesFuncs(
			page_data.TemplateName,
			domain_fm,
			files...,
		)
	}
}

func (w *WebLoader) isTemplate(file_name string) bool {
	return isAnyExt(file_name, w.cfg.TemplatesExtensions)
}

func (w *WebLoader) isResource(file_name string) bool {
	return isAnyExt(file_name, w.cfg.ResourcesExtensions)
}

func (w *WebLoader) getTemplateName(template_file_path string) string {
	rel_path, err := filepath.Rel(w.cfg.DomainDirName, template_file_path)

	if err != nil {
		log.Panicf("%s\n", err.Error())
	}

	file_dir, file_name := filepath.Split(rel_path)
	file_dir = filepath.Clean(file_dir)

	entry_points := []string{w.cfg.PageFileName, w.cfg.LayoutFileName, "index.html"}

	if slices.Contains(entry_points, file_name) {
		return file_dir
	}

	no_ext_filename := strings.TrimSuffix(file_name, filepath.Ext(file_name))

	return filepath.Join(file_dir, no_ext_filename)
}
