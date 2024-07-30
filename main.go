package main

import (
	"allprojects/globalchat"
	"allprojects/imageresizer"

	// "allprojects/urlshortener"
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
    //LOAD  .env file
    //only FOR localhost
//     enverr := godotenv.Load()
//    if enverr != nil {
//        log.Fatalf("Error loading .env file")
//    }
    //
    //
    dbServer := os.Getenv("AZURE_DB_SERVER")
    dbUser := os.Getenv("AZURE_DB_USER")
    dbPassword := os.Getenv("AZURE_DB_PASSWORD")
    dbPort := os.Getenv("AZURE_DB_PORT")
    //dbNameglobalchat for globalchatproject
    dbNameglobalchat := os.Getenv("AZURE_DB_NAME_Global_Chat")
    //dbname for urlshortener project
    // dbNameurl:=os.Getenv("AZURE_DB_NAME_URL")
    //dbconnection for globalchat project
    if dbServer == "" || dbUser == "" || dbPassword == "" || dbPort == "" || dbNameglobalchat == "" {
        log.Fatalf("Database environment variables are not set for globalchat project")
    }
    connectionStringglobalchat := fmt.Sprintf("server=%s;user id=%s;password=%s;port=%s;database=%s;",
        dbServer, dbUser, dbPassword, dbPort, dbNameglobalchat)
    errglobalchat := retryConnect(3, 5*time.Second, func() error {
        return globalchat.InitDB(connectionStringglobalchat)
    })
    if errglobalchat != nil {
        log.Fatalf("Failed to connect to database globalchat project: %v", errglobalchat)
    }
    //dbconnection for urlshortener project
    // if dbServer == "" || dbUser == "" || dbPassword == "" || dbPort == "" || dbNameurl == "" {
    //     log.Fatalf("Database environment variables are not set for urlshortener project")
    // }
    // connectionStringurl := fmt.Sprintf("server=%s;user id=%s;password=%s;port=%s;database=%s;",
    //     dbServer, dbUser, dbPassword, dbPort, dbNameurl)
    // errurl:= retryConnect(3, 5*time.Second, func() error {
    //     return urlshortener.InitDB(connectionStringurl)
    // })
    // if errurl != nil { 
    //     log.Fatalf("Failed to connect to database urlshorterner project: %v", errurl)
    // }
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
    // http.HandleFunc("/imageresizer/download/", imageresizer.DownloadHandler)

    // URL Shortener
    http.HandleFunc("/urlshortener", func(w http.ResponseWriter, r *http.Request) {
        //redirect to website 
        http.Redirect(w, r, "https://shrt.000.pe/", http.StatusSeeOther)
    })
    http.HandleFunc("/urlshortener/", func(w http.ResponseWriter, r *http.Request) {
        //redirect to website
        http.Redirect(w, r, "https://shrt.000.pe//", http.StatusSeeOther)
    })
    

    // http.HandleFunc("/urlshortener", func(w http.ResponseWriter, r *http.Request) {
    //     http.ServeFile(w, r, "urlshortener/static/index.html")
    // })
    // http.HandleFunc("/urlshortener/", func(w http.ResponseWriter, r *http.Request) {
    //     http.ServeFile(w, r, "urlshortener/static/index.html")
    // })
    // http.HandleFunc("/urlshortener/shorten", urlshortener.Shortenurlhandler)

    // Global Chat 
    http.HandleFunc("/globalchat", globalchat.GlobalChatHandler) 
    http.HandleFunc("/globalchat/", globalchat.GlobalChatHandler) 
    http.HandleFunc("/globalchat/send", globalchat.GlobalChatSendHandler) 
    http.HandleFunc("/globalchat/getmessages", globalchat.GlobalChatGetMessagesHandler)

    fmt.Println("Server started at :8080")
    log.Fatal(http.ListenAndServe(":8080", nil))
}