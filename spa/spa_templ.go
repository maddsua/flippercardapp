package spa

import (
	"bytes"
	"fmt"
	"image"
	"io"
	"io/fs"
	"log/slog"
	"regexp"
	"strings"
	"time"

	_ "image/jpeg"
	_ "image/png"
)

func spaLoadIndexTemplateFS(fs fs.FS, name string) ([]byte, error) {

	file, err := fs.Open(name)
	if err != nil {
		return nil, fmt.Errorf("%s not found", name)
	}

	return io.ReadAll(file)
}

type spaIndexProps struct {
	AppDomain string
	AppName   string
}

type spaIndex struct {
	data  []byte
	mtime time.Time
}

func spaGenerateIndex(fs fs.FS, props spaIndexProps) (*spaIndex, error) {

	data, err := spaLoadIndexTemplateFS(fs, "index.html")
	if err != nil {
		return nil, fmt.Errorf("load template: %v", err)
	}

	templProps := map[string]any{
		"og_site_name": "FlipperCard",
	}

	if props.AppDomain != "" {

		ogPreview, err := spaGenerateOgPreview(fs, props.AppDomain)
		if err != nil {
			slog.Warn("SPA: Generate og:preview",
				slog.String("err", err.Error()))
		} else {
			templProps["og_preview"] = ogPreview
		}
	}

	if props.AppName != "" {
		templProps["og_site_name"] = strings.ReplaceAll(props.AppName, `"`, `'`)
	}

	if data, err = MetaTemplateReplace(data, templProps); err != nil {
		return nil, fmt.Errorf("exec template: %v", err)
	}

	page := spaIndex{
		data: data,
	}

	if mt, ok := fs.(modtimer); ok {
		page.mtime = mt.Mtime()
	}

	return &page, nil
}

func spaFindOgPreview(fs fs.FS) (fs.File, error) {

	for _, ext := range []string{"png", "jpg", "jpeg"} {
		if file, _ := fs.Open("og_preview." + ext); file != nil {
			return file, nil
		}
	}

	return nil, fmt.Errorf("preview not found")
}

func spaGenerateOgPreview(fs fs.FS, domain string) (string, error) {

	file, err := spaFindOgPreview(fs)
	if err != nil {
		return "", err
	}

	defer file.Close()

	stat, err := file.Stat()
	if err != nil {
		return "", err
	}

	img, mimetype, err := image.Decode(file)
	if err != nil {
		return "", err
	}

	bounds := img.Bounds()

	return fmt.Sprintf(`
		<meta property="og:image:width" content="%d" />
		<meta property="og:image:height" content="%d" />
		<meta property="og:image:type" content="image/%s" />
		<meta property="og:image" content="https://%s/%s" />
	`, bounds.Dx(), bounds.Dy(), mimetype, domain, stat.Name()), nil
}

func MetaTemplateReplace(data []byte, props map[string]any) ([]byte, error) {

	templRegexp, err := regexp.Compile(`(?i)(\<\!\-{2}\s*\{{2}\s*[^{}]+\s*\}{2}\s*\-{2}\>)|(\{{2}\s*[^{}]+\s*\}{2})`)
	if err != nil {
		return nil, fmt.Errorf("template regexp: %v", err)
	}

	propRegexp, err := regexp.Compile(`(?i)[a-z0-9_]+`)
	if err != nil {
		return nil, fmt.Errorf("prop regexp: %v", err)
	}

	var buff bytes.Buffer

	var lastIdx int

	for _, tokenRange := range templRegexp.FindAllIndex(data, -1) {

		if len(tokenRange) != 2 {
			return nil, fmt.Errorf("invalid regexp match range length (%d)", len(tokenRange))
		}

		tokenStart := tokenRange[0]
		tokenEnd := tokenRange[1]

		if _, err := buff.Write(data[lastIdx:tokenStart]); err != nil {
			return nil, fmt.Errorf("buffer: %v", err)
		}

		lastIdx = tokenEnd
		token := data[tokenStart:tokenEnd]

		tokenProp := propRegexp.Find(token)
		if len(tokenProp) == 0 {
			continue
		}

		value, ok := props[string(tokenProp)]
		if !ok {
			continue
		}

		if _, err := fmt.Fprint(&buff, value); err != nil {
			return nil, fmt.Errorf("buffer: %v", err)
		}
	}

	if lastIdx < len(data)-1 {
		buff.Write(data[lastIdx:])
	}

	return buff.Bytes(), nil
}
