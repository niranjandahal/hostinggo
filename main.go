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
	"time"
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

func retryConnect(attempts int, sleep time.Duration, f func() error) error {
    for i := 0; ; i++ {
        err := f()
        if err == nil {
            return nil
        }

        if i >= (attempts - 1) {
            return err
        }

        time.Sleep(sleep)
    }
}

func main() {

   // only for running on localhost
   //comment this block for deployement

//    err := godotenv.Load()
//    if err != nil {
//        log.Fatalf("Error loading .env file")
//    }
//    dbServer := os.Getenv("AZURE_DB_SERVER")
//    dbUser := os.Getenv("AZURE_DB_USER")
//    dbPassword := os.Getenv("AZURE_DB_PASSWORD")
//    dbPort := os.Getenv("AZURE_DB_PORT")
//    dbName := os.Getenv("AZURE_DB_NAME")
//    if dbServer == "" || dbUser == "" || dbPassword == "" || dbPort == "" || dbName == "" {
//        log.Fatalf("Database environment variables are not set")
//    }
//    connectionString := fmt.Sprintf("server=%s;user id=%s;password=%s;port=%s;database=%s;",
//        dbServer, dbUser, dbPassword, dbPort, dbName)

//    err = retryConnect(3, 5*time.Second, func() error {
//        return globalchat.InitDB(connectionString)
//    })
//    if err != nil {
//        log.Fatalf("Failed to connect to database: %v", err)
//    }

    //
    //
    //
    //uncomment this block for deployement 
    //
    //
    //

    dbServer := os.Getenv("AZURE_DB_SERVER")
    dbUser := os.Getenv("AZURE_DB_USER")
    dbPassword := os.Getenv("AZURE_DB_PASSWORD")
    dbPort := os.Getenv("AZURE_DB_PORT")
    dbName := os.Getenv("AZURE_DB_NAME")

    if dbServer == "" || dbUser == "" || dbPassword == "" || dbPort == "" || dbName == "" {
        log.Fatalf("Database environment variables are not set")
    }

    connectionString := fmt.Sprintf("server=%s;user id=%s;password=%s;port=%s;database=%s;",
        dbServer, dbUser, dbPassword, dbPort, dbName)

    err := retryConnect(3, 5*time.Second, func() error {
        return globalchat.InitDB(connectionString)
    })
    if err != nil {
        log.Fatalf("Failed to connect to database: %v", err)
    }

    //
    //
    //
    //handlers for all projects 
    //
    //
    //
    http.HandleFunc("/", mainHandler)

    // Image Resizer 
    http.HandleFunc("/imageresizer", func(w http.ResponseWriter, r *http.Request) {
        http.ServeFile(w, r, "imageresizer/static/imageresizer.html")
    })
    http.HandleFunc("/imageresizer/upload", imageresizer.UploadHandler)
    http.HandleFunc("/imageresizer/download/", imageresizer.DownloadHandler)

    // URL Shortener 
    http.HandleFunc("/urlshortener", func(w http.ResponseWriter, r *http.Request) {
        http.ServeFile(w, r, "urlshortener/static/urlshortener.html")
    })
    http.HandleFunc("/urlshortener/shorten", urlshortener.ShortenHandler)
    http.HandleFunc("/urlshortener/redirect/", urlshortener.RedirectHandler)

    // Global Chat 
    http.HandleFunc("/globalchat", globalchat.GlobalChatHandler) 
    http.HandleFunc("/globalchat/send", globalchat.GlobalChatSendHandler) 
    http.HandleFunc("/globalchat/getmessages", globalchat.GlobalChatGetMessagesHandler)

    fmt.Println("Server started at :8080")
    log.Fatal(http.ListenAndServe(":8080", nil))
}
