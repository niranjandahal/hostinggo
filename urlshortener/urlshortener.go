package urlshortener

import (
	"html/template"
	"math/rand"
	"net/http"
	"strings"
	"sync"
	"time"
)

var (
	urlMap   = make(map[string]string)
	urlMutex sync.RWMutex
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

func generateShortURL() string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	b := make([]byte, 6)
	for i := range b {
		b[i] = charset[rand.Intn(len(charset))]
	}
	return string(b)
}

func ShortenHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		longURL := r.FormValue("url")
		if longURL == "" {
			http.Error(w, "URL is required", http.StatusBadRequest)
			return
		}

		shortURL := generateShortURL()
		urlMutex.Lock()
		urlMap[shortURL] = longURL
		urlMutex.Unlock()

		data := map[string]string{
			"ShortURL": r.Host + "/urlshortener/redirect/" + shortURL,
		}

		tmpl, err := template.ParseFiles("urlshortener/static/urlshortener.html")
		if err != nil {
			http.Error(w, "Failed to load template", http.StatusInternalServerError)
			return
		}
		tmpl.Execute(w, data)
	} else {
		http.ServeFile(w, r, "urlshortener/static/urlshortener.html")
	}
}

func RedirectHandler(w http.ResponseWriter, r *http.Request) {
	shortURL := strings.TrimPrefix(r.URL.Path, "/urlshortener/redirect/")
	urlMutex.RLock()
	longURL, exists := urlMap[shortURL]
	urlMutex.RUnlock()
	if !exists {
		http.NotFound(w, r)
		return
	}
	http.Redirect(w, r, longURL, http.StatusFound)
}
