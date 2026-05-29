package spa

import (
	"io"
	"io/fs"
	"net/http"
	"path"
	"time"
)

func NewServerSPA(fsys fs.FS, prefix string) http.Handler {
	return http.HandlerFunc(func(wrt http.ResponseWriter, req *http.Request) {

		pathname := req.URL.Path

		if loc, ok := matchRewrite(pathname,
			rewriteIndexSuffix,
			rewriteTrailingHtmlSuffix,
			rewriteTrailingSlash,
		); ok {
			urlRedirect(wrt, loc)
			return
		}

		asset := findWebAsset(fsys, prefix, pathname)
		if asset == nil {
			wrt.WriteHeader(http.StatusNotFound)
			wrt.Write([]byte("asset not found\r\n"))
			return
		}

		defer asset.Close()

		if asset.modtime.IsZero() {
			wrt.Header().Add("Cache-Control", "no-cache")
		}

		http.ServeContent(wrt, req, asset.name, asset.modtime, asset.file.(io.ReadSeeker))
	})
}

func webAssetPath(prefix, name string) string {
	return path.Join(prefix, path.Clean(name))
}

type webAsset struct {
	name    string
	modtime time.Time
	file    fs.File
}

func (asset *webAsset) Close() error {
	if asset == nil || asset.file == nil {
		return nil
	}
	return asset.file.Close()
}

func findWebAsset(fsys fs.FS, prefix, name string) *webAsset {

	assetName := webAssetPath(prefix, name)

	file, _ := fsys.Open(assetName)
	if file == nil {
		// fallback case for SPA routing
		return tryWebAssetPath(fsys, prefix, "/index.html")
	}

	stat, _ := file.Stat()
	if stat == nil {
		return nil
	}

	if stat.IsDir() {
		file.Close()
		return tryWebAssetPath(fsys, prefix, path.Join(name, "/index.html"))
	}

	return &webAsset{assetName, stat.ModTime(), file}
}

func tryWebAssetPath(fsys fs.FS, prefix, name string) *webAsset {
	return tryWebAsset(fsys, webAssetPath(prefix, name))
}

func tryWebAsset(fsys fs.FS, asset string) *webAsset {

	file, _ := fsys.Open(asset)
	if file == nil {
		return nil
	}

	stat, _ := file.Stat()
	if stat == nil || stat.IsDir() {
		return nil
	}

	return &webAsset{asset, stat.ModTime(), file}
}

type rewriteRuleFn func(location string) (string, bool)

func matchRewrite(location string, rules ...rewriteRuleFn) (string, bool) {
	for _, rule := range rules {
		if loc, ok := rule(location); ok {
			return loc, ok
		}
	}
	return "", false
}

func urlRedirect(wrt http.ResponseWriter, location string) {
	wrt.Header().Set("Location", location)
	wrt.WriteHeader(http.StatusMovedPermanently)
}
