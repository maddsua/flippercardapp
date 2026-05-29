package spa

import "strings"

type Rewrite interface {
	Rewrite(location string) (string, bool)
}

type IndexSuffixRewrite struct{}

func (rewrite IndexSuffixRewrite) Rewrite(location string) (string, bool) {

	clean, found := strings.CutSuffix(location, "index.html")
	if !found {
		return "", false
	}

	if clean == "" {
		clean = "/"
	}

	return clean, true
}

type TrailingSuffixRewrite struct{}

func (rewrite TrailingSuffixRewrite) Rewrite(location string) (string, bool) {
	if clean, found := strings.CutSuffix(location, ".html"); found {
		return clean, found
	}
	return strings.CutSuffix(location, ".htm")
}

type TrailingSlashRewrite struct{}

func (rewrite TrailingSlashRewrite) Rewrite(location string) (string, bool) {
	if location == "/" {
		return "", false
	}

	return strings.CutSuffix(location, "/")
}
