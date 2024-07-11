package urlshortener

import (
	"net/http"
)

func ShortenHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "urlshortener/static/urlshortener.html")
}

func RedirectHandler(w http.ResponseWriter, r *http.Request) {
	// Implement URL redirection logic here
	http.Redirect(w, r, "http://example.com", http.StatusFound)
}
