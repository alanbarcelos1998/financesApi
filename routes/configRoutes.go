package routes

import "github.com/gorilla/mux"

type ResponseErr struct {
	Erro string `json:"erro"`
}

func ConfigRoute(router *mux.Router) {
	// Fixed Expenses Routes
	router.HandleFunc("/fixedexpensesdate/{dateregister}", filterFixedDate).Methods("GET")
	router.HandleFunc("/fixedexpenses", listFixedExpenses).Methods("GET")
	router.HandleFunc("/addfixedexpense", addFixedExpense).Methods("POST")
	router.HandleFunc("/fixedexpenses/{idfixed}", alterFixedExpense).Methods("PUT")
	router.HandleFunc("/fixedexpenses/{idfixed}", deleteFixed).Methods("DELETE")

	// Variable Expenses Routes
	router.HandleFunc("/variableexpensesdate/{dateregister}", filterVariableDate).Methods("GET")
	router.HandleFunc("/variableexpenses", listVariableExpenses).Methods("GET")
	router.HandleFunc("/addvariableexpense", addVariableExpense).Methods("POST")
	router.HandleFunc("/variableexpenses/{idvariable}", alterVariableExpense).Methods("PUT")
	router.HandleFunc("/variableexpenses/{idvariable}", deleteVariable).Methods("DELETE")
}
