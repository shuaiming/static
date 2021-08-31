package static

import (
	"net/http"
	"strings"
)

// Static serve static files
type Static struct {
	Prefix string
	Dir    http.FileSystem
	Index  bool
}

// New Static
func New(prefix string, fs http.FileSystem, index bool) *Static {

	return &Static{
		Prefix: prefix, // statics url prefix
		Dir:    fs,
		Index:  index,
	}
}

// ServeHTTP implement pod.Handler
func (s *Static) ServeHTTP(
	w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {

	if !strings.HasPrefix(r.URL.Path, s.Prefix) {
		next(w, r)
		return
	}

	// do not list dir when needed
	if !s.Index && strings.HasSuffix(r.URL.Path, "/") {
		http.Error(w, "403 Forbidden", http.StatusForbidden)
		return
	}

	// fix http.StripPrefix dirList() bug
	if r.URL.Path == s.Prefix && !strings.HasSuffix(r.URL.Path, "/") {
		http.Redirect(w, r, r.URL.Path+"/", http.StatusMovedPermanently)
		return
	}

	http.StripPrefix(s.Prefix, http.FileServer(s.Dir)).ServeHTTP(w, r)
}
