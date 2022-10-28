package main

import (
	"automation-project/db"
	"automation-project/handler"
	"automation-project/model"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/rs/cors"
)

func main() {
	//Init DB
	db, err := db.InitDB()
	if err != nil {
		log.Fatal(err)
	}

	//Init model
	mapModel := &model.MapCRUD{
		DB: db,
	}

	//Init handler
	mapHandler := &handler.MapHandler{
		MapModel: mapModel,
	}

	//Init router
	r := mux.NewRouter()

	//routes
	r.HandleFunc("/", mapHandler.Create).Methods("POST")
	r.HandleFunc("/", mapHandler.List).Methods("GET")
	r.HandleFunc("/get", mapHandler.GetByKey).Methods("GET")
	r.HandleFunc("/", mapHandler.Update).Methods("PUT")
	r.HandleFunc("/", mapHandler.Delete).Methods("DELETE")

	fmt.Println("HTTP server running on http://127.0.0.1:8080")

	handler := cors.AllowAll().Handler(r)
	err = http.ListenAndServe(":8080", handler)
	if err != nil {
		log.Fatal(err)
	}
}
