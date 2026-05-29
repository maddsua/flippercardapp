package spa

import (
	"bytes"
	"fmt"
	"html/template"
	"image"
	"io"
	"io/fs"
	"log/slog"
	"os"
	"time"

	_ "image/jpeg"
	_ "image/png"
)

func loadTemplateFileFs(fs fs.FS, name string) (*template.Template, error) {

	file, err := fs.Open(name)
	if err != nil {
		return nil, fmt.Errorf("index.html not found")
	}

	page, err := io.ReadAll(file)
	if err != nil {
		return nil, err
	}

	return template.New(name).Parse(string(page))
}

type pageTemplateData struct {
	data  []byte
	mtime time.Time
}

func spaIndexTemplate(fs fs.FS) (*pageTemplateData, error) {

	templ, err := loadTemplateFileFs(fs, "index.html")
	if err != nil {
		return nil, err
	}

	props := map[string]any{
		"app_domain": "localhost",
		"app_name":   "FlipperCard",
	}

	if val := os.Getenv("APP_DOMAIN"); val != "" {

		props["app_domain"] = val

		ogPreview, err := spaGenerateOgPreview(fs, val)
		if err != nil {
			slog.Warn("SPA: Generate og:preview",
				slog.String("err", err.Error()))
		} else {
			props["og_preview"] = ogPreview
		}
	}

	if val := os.Getenv("APP_NAME"); val != "" {
		props["app_name"] = val
	}

	var buff bytes.Buffer
	if err := templ.Execute(&buff, props); err != nil {
		return nil, err
	}

	page := pageTemplateData{
		data: buff.Bytes(),
	}

	if mt, ok := fs.(modtimer); ok {
		page.mtime = mt.Mtime()
	}

	return &page, nil
}

func spaGenerateOgPreview(fs fs.FS, domain string) (template.HTML, error) {

	for _, ext := range []string{"png", "jpg", "jpeg"} {
		if file, _ := fs.Open("og_preview." + ext); file != nil {
			return spaGenerateOgPreviewFromFile(file, domain)
		}
	}

	return "", fmt.Errorf("preview not found")
}

func spaGenerateOgPreviewFromFile(file fs.File, domain string) (template.HTML, error) {

	defer file.Close()

	stat, err := file.Stat()
	if err != nil {
		return "", err
	}

	img, mimetype, err := image.Decode(file)
	if err != nil {
		return "", err
	}

	return spaGenerateOgPreviewTemplate(domain, stat.Name(), "image/"+mimetype, img.Bounds().Dx(), img.Bounds().Dy())
}

func spaGenerateOgPreviewTemplate(domain, name, mimetype string, width, height int) (template.HTML, error) {

	templ, err := template.New("og-preview").Parse(`
		<meta property="og:image:width" content="{{.width}}" />
		<meta property="og:image:height" content="{{.height}}" />
		<meta property="og:image:type" content="{{.type}}" />
		<meta property="og:image" content="https://{{.domain}}/{{.name}}" />
	`)

	if err != nil {
		return "", err
	}

	var buff bytes.Buffer

	if err := templ.Execute(&buff, map[string]any{
		"width":  width,
		"height": height,
		"domain": domain,
		"name":   name,
		"type":   mimetype,
	}); err != nil {
		return "", err
	}

	return template.HTML(buff.String()), nil
}
