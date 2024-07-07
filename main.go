package main

import (
    "log"
    "net/http"
    "hostedgo/database"
    "hostedgo/handlers"

    "github.com/gorilla/mux"
)

func main() {
    var err error
    database.DB, err = database.InitDB()
    if err != nil {
        log.Fatal(err)
    }

    router := mux.NewRouter()

    router.HandleFunc("/items", handlers.CreateItem).Methods("POST")
    router.HandleFunc("/items", handlers.GetItems).Methods("GET")
    router.HandleFunc("/items/{id}", handlers.GetItem).Methods("GET")
    router.HandleFunc("/items/{id}", handlers.UpdateItem).Methods("PUT")
    router.HandleFunc("/items/{id}", handlers.DeleteItem).Methods("DELETE")

    log.Fatal(http.ListenAndServe(":8000", router))
}
