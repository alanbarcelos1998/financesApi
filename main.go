package main

import (
	db "financasApi/db"
	route "financasApi/routes"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func jsonMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		next.ServeHTTP(w, r)
	})
}

func configServer() {
	router := mux.NewRouter().StrictSlash(true)
	router.Use(jsonMiddleware)
	route.ConfigRoute(router)

	fmt.Println("Server is running port 8785")
	log.Fatal(http.ListenAndServe(":8785", router))
}

func main() {
	db.Connection()
	configServer()
}
