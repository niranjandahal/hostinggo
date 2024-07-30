package urlshortener

import (
	"database/sql"
	"fmt"
	"html/template"
	"math/rand"
	"net/http"

	_ "github.com/microsoft/go-mssqldb"
)

var db *sql.DB

func InitDB(dataSourceName string) error {
	var err error
	db, err = sql.Open("sqlserver", dataSourceName)
	if err != nil {
		return fmt.Errorf("failed to connect to database urlshortener: %v", err)
	}

	err = db.Ping()
	if err != nil {
		return fmt.Errorf("failed to ping database urlshortener: %v", err)
	}

	fmt.Println("Connected to database of urlshortener")
	return nil
}

func shrt_url_code_gen() string {
	var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")
	b := make([]rune, 5)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}

func Shortenurlhandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	err := r.ParseForm()
	if err != nil {
		http.Error(w, "Failed to parse form", http.StatusInternalServerError)
		return
	}

	longURL := r.FormValue("url")
	customCode := r.FormValue("customCode")

	fmt.Println("Long URL: ", longURL)
	fmt.Println("Custom Code: ", customCode)

	if longURL == "" {
		http.Error(w, "URL is required", http.StatusBadRequest)
		return
	}

	var shortURL string
	var urlID int64
	var customURL string
	var customcodeerrmsg string

	err = db.QueryRow("SELECT id FROM urls WHERE actual_url = @p1", longURL).Scan(&urlID)
	if err == sql.ErrNoRows {
		shortURL = shrt_url_code_gen()
		var newID int64
		query := "INSERT INTO urls (actual_url, shrt_url_code) OUTPUT INSERTED.id VALUES (@p1, @p2)"
		err = db.QueryRow(query, longURL, shortURL).Scan(&newID)
		if err != nil {
			http.Error(w, "Failed to insert URL into database", http.StatusInternalServerError)
			fmt.Println("Error inserting URL:", err)
			return
		}
		urlID = newID
	} else if err == nil {
		err = db.QueryRow("SELECT shrt_url_code FROM urls WHERE id = @p1", urlID).Scan(&shortURL)
		if err != nil {
			http.Error(w, "Failed to retrieve existing short URL code", http.StatusInternalServerError)
			fmt.Println("Error retrieving short URL code:", err)
			return
		}
	} else {
		http.Error(w, "Failed to query URLs", http.StatusInternalServerError)
		fmt.Println("Error querying URLs:", err)
		return
	}

	if customCode != "" {
		var existingURLID int64
		customURL = "shrt.000.pe/" + customCode

		err = db.QueryRow("SELECT url_id FROM custom_codes WHERE custom_code = @p1", customCode).Scan(&existingURLID)
		if err == nil && existingURLID != urlID {
			customcodeerrmsg = "Custom code already in use. Try a different one."
			customURL = customCode + " is already in use."
		} else if err != nil && err != sql.ErrNoRows {
			http.Error(w, "Failed to query custom codes", http.StatusInternalServerError)
			fmt.Println("Error querying custom codes:", err)
			return
		} else {
			_, err = db.Exec("INSERT INTO custom_codes (url_id, custom_code) VALUES (@p1, @p2)", urlID, customCode)
			if err != nil {
				http.Error(w, "Failed to insert custom code into database", http.StatusInternalServerError)
				fmt.Println("Error inserting custom code:", err)
				return
			}
		}
	}

	data := map[string]string{
		"ShortURL":        "shrt.000.pe/" + shortURL,
		"CustomURL":       customURL,
		"CustomCodeError": customcodeerrmsg,
	}

	tmpl, err := template.ParseFiles("urlshortener/static/urlshortener.html")
	if err != nil {
		http.Error(w, "Failed to load template", http.StatusInternalServerError)
		fmt.Println("Error loading template:", err)
		return
	}
	tmpl.Execute(w, data)
}

