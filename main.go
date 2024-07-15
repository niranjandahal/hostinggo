package main

import (
	"allprojects/globalchat"
	"allprojects/imageresizer"
	"allprojects/urlshortener"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"github.com/joho/godotenv"
)

type Project struct {
    Name string
    URL  string
}

var projects = []Project{
    {"Global Chat", "/globalchat"},
    {"Image Resizer", "/imageresizer"},
    {"URL Shortener", "/urlshortener"},
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
    // Load environment variables from .env file
    err := godotenv.Load()
    if err != nil {
        log.Fatalf("Error loading .env file")
    }

    // Load environment variables
    dbServer := os.Getenv("AZURE_DB_SERVER")
    dbUser := os.Getenv("AZURE_DB_USER")
    dbPassword := os.Getenv("AZURE_DB_PASSWORD")
    dbPort := os.Getenv("AZURE_DB_PORT")
    dbName := os.Getenv("AZURE_DB_NAME")

    if dbServer == "" || dbUser == "" || dbPassword == "" || dbPort == "" || dbName == "" {
        log.Fatalf("Database environment variables are not set")
    }

    // Initialize database connection
    connectionString := fmt.Sprintf("server=%s;user id=%s;password=%s;port=%s;database=%s;",
        dbServer, dbUser, dbPassword, dbPort, dbName)
    globalchat.InitDB(connectionString)

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
    http.HandleFunc("/globalchat/getmessages", globalchat.GlobalChatGetMessagesHandler) // Route for getting messages

    // Start server
    fmt.Println("Server started at :8080")
    log.Fatal(http.ListenAndServe(":8080", nil))
}
