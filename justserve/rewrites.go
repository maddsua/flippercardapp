package justserve

import "strings"

func rewriteIndexSuffix(location string) (string, bool) {

	clean, found := strings.CutSuffix(location, "index.html")
	if !found {
		return "", false
	}

	if clean == "" {
		clean = "/"
	}

	return clean, true
}

func rewriteTrailingHtmlSuffix(location string) (string, bool) {
	return strings.CutSuffix(location, ".html")
}

func rewriteTrailingSlash(location string) (string, bool) {

	if location == "/" {
		return "", false
	}

	return strings.CutSuffix(location, "/")
}
