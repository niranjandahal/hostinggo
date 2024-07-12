package main

import (
	"allprojects/globalchat"
	"allprojects/imageresizer"
	"allprojects/urlshortener"
	"fmt"
	"html/template"
	"log"
	"net/http"
)

type Project struct {
	Name string
	URL  string
}

var projects = []Project{
	{"Image Resizer", "/imageresizer"},
	{"URL Shortener", "/urlshortener"},
	{"Global Chat", "/globalchat"},
}

func mainHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	tmpl, err := template.ParseFiles("templates/index.html")
	if err != nil {
		http.Error(w, "Failed to load template", http.StatusInternalServerError)
		return
	}

	tmpl.Execute(w, projects)
}

func main() {
	// Initialize database
	globalchat.InitDB("root:@tcp(localhost:3306)/globalchat")

	// Handle routes
	http.HandleFunc("/", mainHandler)

	// Image Resizer routes
	http.HandleFunc("/imageresizer", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "imageresizer/static/imageresizer.html")
	})
	http.HandleFunc("/imageresizer/upload", imageresizer.UploadHandler)
	http.HandleFunc("/imageresizer/download/", imageresizer.DownloadHandler)

	// URL Shortener routes
	http.HandleFunc("/urlshortener", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "urlshortener/static/urlshortener.html")
	})
	http.HandleFunc("/urlshortener/shorten", urlshortener.ShortenHandler)
	http.HandleFunc("/urlshortener/redirect/", urlshortener.RedirectHandler)

	// Global Chat routes
	http.HandleFunc("/globalchat", globalchat.GlobalChatHandler) // Route for fetching messages
	http.HandleFunc("/globalchat/send", globalchat.GlobalChatSendHandler) // Route for sending messages
    http.HandleFunc("/globalchat/getmessages", globalchat.GlobalChatGetMessagesHandler)


	// Start server
	fmt.Println("Server started at :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
