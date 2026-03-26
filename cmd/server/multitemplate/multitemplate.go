package multitemplate

import (
	"io/fs"
	"log"
	"path/filepath"
	"slices"
	"strings"

	"github.com/gin-contrib/multitemplate"
	"github.com/gin-gonic/gin"
	"github.com/yogusita/to-adhdo/cmd/server/funcmap"
)

type WebLoaderConfig struct {
	PageFileName   string
	LayoutFileName string

	DomainDirName    string
	WebSourceDirName string
	PagesDirName     string

	TemplatesExtensions []string
	ResourcesExtensions []string
}

var cfg = WebLoaderConfig{
	LayoutFileName: "layout.html",

	DomainDirName:    "domain",
	WebSourceDirName: "web",
	PagesDirName:     "pages",

	TemplatesExtensions: []string{
		".html",
	},
	ResourcesExtensions: []string{
		".js",
		".css",
		".svg",
		".img",
		".png",
		".jpg",
	},
}

func DefaultWebConfig() WebLoaderConfig {
	return cfg
}

func matchExt(file_name string, extensions []string) bool {
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
	DomainName string
	Components []TemplateData
	Pages      []PageTemplateData
}

type WebLoader struct {
	fs   fs.FS
	cfg  WebLoaderConfig
	data []DomainTemplatesData
}

func CreateDefaultWebLoader(fsys fs.FS) *WebLoader {
	return &WebLoader{
		fs:  fsys,
		cfg: DefaultWebConfig(),
	}
}

func (w *WebLoader) LoadTemplates() (multitemplate.Renderer, error) {
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

		w.loadDomainFiles(&r, fm, domain_data)
	}

	return r, nil
}

func (w *WebLoader) RouteResources(router *gin.Engine) error {
	domains, err := fs.ReadDir(w.fs, w.cfg.DomainDirName)

	if err != nil {
		return err
	}

	for _, domain_entry := range domains {
		if !domain_entry.IsDir() {
			continue
		}

		domain_name := domain_entry.Name()
		domain_dir_path := filepath.Join(w.cfg.DomainDirName, domain_name)

		err := fs.WalkDir(w.fs, domain_dir_path, func(path string, d fs.DirEntry, err error) error {
			if err != nil {
				return err
			}

			if d.IsDir() || !w.isResource(d.Name()) {
				return nil
			}

			router.StaticFile(path, path)

			return nil
		})

		if err != nil {
			return err
		}
	}

	return nil
}

func (w *WebLoader) webRootDirPath(domain_name string) string {
	return filepath.Join(w.cfg.DomainDirName, domain_name, w.cfg.WebSourceDirName)
}

func (w *WebLoader) pagesDirPath(domain_name string) string {
	return filepath.Join(w.webRootDirPath(domain_name), w.cfg.PagesDirName)
}

func (w *WebLoader) defaultLayoutPath(domain_name string) string {
	return filepath.Join(w.webRootDirPath(domain_name), w.cfg.LayoutFileName)
}

func (w *WebLoader) pageLayoutPath(domain_name, page_name string) string {
	p := filepath.Join(w.pagesDirPath(domain_name), page_name, cfg.LayoutFileName)

	if _, err := fs.Stat(w.fs, p); err != nil {
		return w.defaultLayoutPath(domain_name)
	}

	return p
}

func (w *WebLoader) isPagePath(domain_name, page_path string) bool {
	input_individual_page_dir, page_index_file := filepath.Split(page_path)
	input_page_dir, _ := filepath.Split(input_individual_page_dir)

	if input_page_dir != w.pagesDirPath(domain_name) {
		return false
	}

	if !w.isTemplate(page_index_file) {
		return false
	}

	if strings.TrimSuffix(page_index_file, filepath.Ext(page_index_file)) != "index" {
		return false
	}

	return true
}

func (w *WebLoader) getDomainData(domain_name string) (DomainTemplatesData, error) {
	data := DomainTemplatesData{
		DomainName: domain_name,
	}

	err := fs.WalkDir(w.fs, w.webRootDirPath(domain_name), func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if d.IsDir() {
			return nil
		}

		if !w.isTemplate(d.Name()) {
			return nil
		}

		if strings.EqualFold(d.Name(), w.cfg.LayoutFileName) {
			return nil
		}

		template_name := w.getTemplateName(path)

		if w.isPagePath(domain_name, path) {
			data.Pages = append(data.Pages, PageTemplateData{
				TemplateName: template_name,
				FilePath:     path,
				LayoutPath:   w.pageLayoutPath(domain_name, filepath.Base(template_name)),
			})

		} else {
			data.Components = append(data.Components, TemplateData{
				TemplateName: template_name,
				FilePath:     path,
			})
		}

		return nil
	})

	if err != nil {
		return data, err
	}

	return data, nil
}

func (w *WebLoader) loadDomainFiles(r *multitemplate.Renderer, fm *funcmap.TemplateFuncMaps, domain_data DomainTemplatesData) {
	domain_fm := fm.FuncMap(domain_data.DomainName)

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
	return matchExt(file_name, w.cfg.TemplatesExtensions)
}

func (w *WebLoader) isResource(file_name string) bool {
	return matchExt(file_name, w.cfg.ResourcesExtensions)
}

func (w *WebLoader) getTemplateName(template_file_path string) string {
	rel_path, err := filepath.Rel(w.cfg.DomainDirName, template_file_path)

	if err != nil {
		log.Panicf("%s\n", err.Error())
	}

	file_dir, file_name := filepath.Split(rel_path)
	file_dir = filepath.Clean(file_dir)

	entry_points := []string{
		w.cfg.PageFileName,
		w.cfg.LayoutFileName,
		"index.html",
	}

	if slices.Contains(entry_points, file_name) {
		return file_dir
	}

	no_ext_filename := strings.TrimSuffix(file_name, filepath.Ext(file_name))

	return filepath.Join(file_dir, no_ext_filename)
}
