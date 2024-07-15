package globalchat

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"

	_ "github.com/microsoft/go-mssqldb"
)

var db *sql.DB

func InitDB(dataSourceName string) error {
    var err error
    db, err = sql.Open("sqlserver", dataSourceName)
    if err != nil {
        return fmt.Errorf("failed to connect to database: %v", err)
    }

    err = db.Ping()
    if err != nil {
        return fmt.Errorf("failed to ping database: %v", err)
    }

    fmt.Print("Connected to database\n")
    return nil
}

type Message struct {
    ID        int
    Content   string
    Timestamp string
}

func fetchMessages() ([]Message, error) {
    rows, err := db.Query("SELECT id, content, timestamp FROM messages ORDER BY timestamp ASC")
    if err != nil {
        log.Printf("Failed to fetch messages from database: %v", err)
        return nil, err
    }
    defer rows.Close()

    var messages []Message

    for rows.Next() {
        var msg Message
        err := rows.Scan(&msg.ID, &msg.Content, &msg.Timestamp)
        if err != nil {
            log.Printf("Error scanning message row: %v", err)
            return nil, err
        }
        messages = append(messages, msg)
    }
    if err := rows.Err(); err != nil {
        log.Printf("Error iterating over message rows: %v", err)
        return nil, err
    }

    log.Printf("Fetched messages: %v", messages) 

    return messages, nil
}

func GlobalChatHandler(w http.ResponseWriter, r *http.Request) {
    messages, err := fetchMessages()
    if err != nil {
        http.Error(w, "Failed to fetch messages from database", http.StatusInternalServerError)
        return
    }

    data := struct {
        Messages []Message
    }{
        Messages: messages,
    }

    renderTemplate(w, "globalchat/static/globalchat.html", data)
}

func GlobalChatGetMessagesHandler(w http.ResponseWriter, r *http.Request) {
    messages, err := fetchMessages()
    if err != nil {
        http.Error(w, "Failed to fetch messages from database", http.StatusInternalServerError)
        return
    }

    jsonResponse, err := json.Marshal(struct {
        Messages []Message `json:"messages"`
    }{Messages: messages})
    if err != nil {
        http.Error(w, "Failed to serialize messages to JSON", http.StatusInternalServerError)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    w.Write(jsonResponse)
}


func GlobalChatSendHandler(w http.ResponseWriter, r *http.Request) {
    fmt.Println("GlobalChatSendHandler function called")

    if r.Method != http.MethodPost {
        http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
        return
    }

    err := r.ParseForm()
    if err != nil {
        http.Error(w, "Failed to parse form", http.StatusInternalServerError)
        return
    }

    message := r.Form.Get("message")
    if message == "" {
        http.Error(w, "Message cannot be empty", http.StatusBadRequest)
        return
    }

    fmt.Printf("Inserting message: %s\n", message)

    result, err := db.Exec("INSERT INTO messages (content) VALUES (@p1)", message)
    // result, err := db.Exec("INSERT INTO messages (content) VALUES (?)", message)
    if err != nil {
        log.Printf("Failed to insert message into database: %v", err)
        http.Error(w, "Failed to insert message into database", http.StatusInternalServerError)
        return
    }

    rowsAffected, err := result.RowsAffected()
    if err != nil {
        log.Printf("Failed to get affected rows: %v", err)
        http.Error(w, "Failed to get affected rows", http.StatusInternalServerError)
        return
    }

    fmt.Printf("Message sent: %s, Rows affected: %d\n", message, rowsAffected)

    http.Redirect(w, r, "/globalchat", http.StatusSeeOther)
}

func renderTemplate(w http.ResponseWriter, tmpl string, data interface{}) {
    t, err := template.ParseFiles(tmpl)
    if err != nil {
        http.Error(w, "Failed to load template", http.StatusInternalServerError)
        return
    }

    err = t.Execute(w, data)
    if err != nil {
        http.Error(w, "Failed to render template", http.StatusInternalServerError)
        return
    }
}
