package spa

import (
	"bytes"
	"fmt"
	"io"
	"io/fs"
	"os"
	"text/template"
	"time"
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

type spaTempalteVars struct {
	AppDomain string
}

func spaIndexTemplate(fs fs.FS) (*pageTemplateData, error) {

	templ, err := loadTemplateFileFs(fs, "index.html")
	if err != nil {
		return nil, err
	}

	vars := spaTempalteVars{
		AppDomain: "localhost",
	}

	if val := os.Getenv("APP_DOMAIN"); val != "" {
		vars.AppDomain = val
	}

	var buff bytes.Buffer
	if err := templ.Execute(&buff, vars); err != nil {
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
