package spa

import (
	"bytes"
	"io"
	"io/fs"
	"log/slog"
	"net/http"
	"path"
	"strings"
	"time"
)

func NewServerSPA(fs fs.FS) http.Handler {

	index, err := spaIndexTemplate(fs)
	if err != nil {
		slog.Warn("SPA: Load index template",
			slog.String("err", err.Error()))
	}

	rewrites := []Rewrite{
		IndexSuffixRewrite{},
		TrailingSuffixRewrite{},
		TrailingSlashRewrite{},
	}

	return http.HandlerFunc(func(wrt http.ResponseWriter, req *http.Request) {

		reqPath := req.URL.Path

		for _, rule := range rewrites {
			if loc, ok := rule.Rewrite(reqPath); ok {
				wrt.Header().Set("Location", loc)
				wrt.WriteHeader(http.StatusMovedPermanently)
				return
			}
		}

		if index != nil && (reqPath == "" || reqPath == "/") {
			serveIndex(wrt, req, reqPath, index)
			return
		}

		if file, _ := fs.Open(reqPath); file != nil {

			defer file.Close()

			if stat, _ := file.Stat(); stat == nil || !stat.Mode().IsRegular() {
				serveIndex(wrt, req, reqPath, index)
				return
			}

			serveFile(wrt, req, reqPath, file)
			return
		}

		if isPageRequest(req) {
			serveIndex(wrt, req, reqPath, index)
			return
		}

		serve404(wrt)
	})
}

func serveFile(wrt http.ResponseWriter, req *http.Request, assetPath string, file fs.File) {

	var modtime time.Time
	if stat, _ := file.Stat(); stat != nil {
		modtime = stat.ModTime()
	}

	reader, ok := file.(io.ReadSeeker)
	if !ok {
		wrt.WriteHeader(http.StatusInternalServerError)
		return
	}

	setCacheControl(wrt)

	http.ServeContent(wrt, req, assetPath, modtime, reader)
}

func serveIndex(wrt http.ResponseWriter, req *http.Request, assetPath string, templ *pageTemplateData) {

	if templ == nil {
		serve404(wrt)
		return
	}

	setCacheControl(wrt)
	http.ServeContent(wrt, req, assetPath, templ.mtime, bytes.NewReader(templ.data))
}

func serve404(wrt http.ResponseWriter) {
	wrt.Header().Set("Cache-Control", "no-cache")
	wrt.WriteHeader(http.StatusNotFound)
	wrt.Write([]byte("asset not found\r\n"))
}

func setCacheControl(wrt http.ResponseWriter) {
	wrt.Header().Set("Cache-Control", "max-age=3600, must-revalidate")
}

func isPageRequest(req *http.Request) bool {

	if path.Ext(req.URL.Path) == "" {
		return true
	}

	for mimetype := range strings.SplitSeq(req.Header.Get("Accept"), ",") {

		mimetype, _, _ = strings.Cut(mimetype, ";")
		_, format, _ := strings.Cut(mimetype, "/")

		if strings.Contains(format, "html") || strings.Contains(format, "xml") {
			return true
		}
	}

	return false
}
