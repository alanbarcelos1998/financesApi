package routes

import "github.com/gorilla/mux"

func ConfigRoute(router *mux.Router) {
	// router.HandleFunc("/", test)
	router.HandleFunc("/fixedexpenses", listFixedExpenses).Methods("GET")
	router.HandleFunc("/addfixedexpense", addFixedExpense).Methods("POST")
	router.HandleFunc("/fixedexpenses/{idfixed}", alterFixedExpense).Methods("PUT")
	router.HandleFunc("/fixedexpenses/{idfixed}", deleteFixed).Methods("DELETE")
	router.HandleFunc("/fixedexpensesdate/{dateregister}", filterDate).Methods("GET")
}
