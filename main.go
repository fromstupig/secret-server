package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/smapig/secret-server/controllers"
)

func main() {
	router := mux.NewRouter()

	router.HandleFunc("/v1/secret", controllers.AddSecret).Methods("POST")
	router.HandleFunc("/v1/secret/{hash}", controllers.GetSecret).Methods("GET")

	sh := http.StripPrefix("/swaggerui/", http.FileServer(http.Dir("./swaggerui/")))
	router.PathPrefix("/swaggerui/").Handler(sh)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	fmt.Println("Maigc happen on ", port)

	err := http.ListenAndServe(":"+port, router)
	if err != nil {
		log.Fatal(err)
	}
}
