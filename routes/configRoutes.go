package routes

import "github.com/gorilla/mux"

type ResponseErr struct {
	Erro string `json:"erro"`
}

func ConfigRoute(router *mux.Router) {
	//Expenses Routes
	router.HandleFunc("/expense/{date}", filterDate).Methods("GET")
	router.HandleFunc("/expense", listExpenses).Methods("GET")
	router.HandleFunc("/expense", cadExpense).Methods("POST")
	router.HandleFunc("/expense/{id}", alterExpense).Methods("PUT")
	router.HandleFunc("/expense/{id}", delete).Methods("DELETE")

}
