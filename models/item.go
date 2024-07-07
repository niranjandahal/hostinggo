package handlers

import (
	"encoding/json"
	"hostedgo/database"
	"net/http"

	"github.com/gorilla/mux"

	"github.com/niranjandahal/hostinggo/models"
)

func CreateItem(w http.ResponseWriter, r *http.Request) {
    var item models.Item
    json.NewDecoder(r.Body).Decode(&item)
    database.DB.Create(&item)
    json.NewEncoder(w).Encode(item)
}

func GetItems(w http.ResponseWriter, r *http.Request) {
    var items []models.Item
    database.DB.Find(&items)
    json.NewEncoder(w).Encode(items)
}

func GetItem(w http.ResponseWriter, r *http.Request) {
    params := mux.Vars(r)
    var item models.Item
    database.DB.First(&item, params["id"])
    json.NewEncoder(w).Encode(item)
}

func UpdateItem(w http.ResponseWriter, r *http.Request) {
    params := mux.Vars(r)
    var item models.Item
    database.DB.First(&item, params["id"])
    json.NewDecoder(r.Body).Decode(&item)
    database.DB.Save(&item)
    json.NewEncoder(w).Encode(item)
}

func DeleteItem(w http.ResponseWriter, r *http.Request) {
    params := mux.Vars(r)
    var item models.Item
    database.DB.Delete(&item, params["id"])
    json.NewEncoder(w).Encode("Item deleted")
}
