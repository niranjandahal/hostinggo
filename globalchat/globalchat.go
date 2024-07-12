package globalchat

import (
	"database/sql"
	"encoding/json"
	"html/template"
	"log"
	"net/http"

	_ "github.com/go-sql-driver/mysql" // MySQL driver
)

var db *sql.DB

func InitDB(dataSourceName string) {
	var err error
	db, err = sql.Open("mysql", dataSourceName)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
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

    log.Printf("Fetched messages: %v", messages) // Log fetched messages

    return messages, nil
}
func GlobalChatHandler(w http.ResponseWriter, r *http.Request) {
    messages, err := fetchMessages()
    if err != nil {
        http.Error(w, "Failed to fetch messages from database", http.StatusInternalServerError)
        return
    }

    // Prepare data to pass into the template
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

    // Prepare JSON response
    jsonResponse, err := json.Marshal(struct {
        Messages []Message `json:"messages"`
    }{Messages: messages})
    if err != nil {
        http.Error(w, "Failed to serialize messages to JSON", http.StatusInternalServerError)
        return
    }

    // Set content type and write JSON response
    w.Header().Set("Content-Type", "application/json")
    w.Write(jsonResponse)
}




func GlobalChatSendHandler(w http.ResponseWriter, r *http.Request) {
	
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

	// Insert message into database
	_, err = db.Exec("INSERT INTO messages (content) VALUES (?)", message)
	if err != nil {
		http.Error(w, "Failed to insert message into database", http.StatusInternalServerError)
		return
	}

	// Redirect back to chat page
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
