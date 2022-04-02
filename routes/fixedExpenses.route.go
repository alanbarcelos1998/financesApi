package routes

import (
	"encoding/json"
	db2 "financasApi/db"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

type FixedExpenses struct {
	Idfixed      int     `json:"id"`
	NameExpense  string  `json:"name"`
	ValueExpense float64 `json:"value"`
	DueDate      string  `json:"duedate"`
	PayDate      string  `json:"paydate"`
}

func listFixedExpenses(w http.ResponseWriter, r *http.Request) {
	registers, errSelect := db2.Db.Query("SELECT * FROM fixed_expenses")

	if errSelect != nil {
		log.Println("fixed_expenses: " + errSelect.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	var fixedExpenses []FixedExpenses = make([]FixedExpenses, 0)
	for registers.Next() {
		var fixedExpense FixedExpenses
		errScan := registers.Scan(&fixedExpense.Idfixed, &fixedExpense.NameExpense, &fixedExpense.ValueExpense, &fixedExpense.DueDate, &fixedExpense.PayDate)
		if errScan != nil {
			log.Println("FixedExepenses: errScan: " + errScan.Error())
			continue
		}

		fixedExpenses = append(fixedExpenses, fixedExpense)
	}

	errCloseRegisters := registers.Close()

	if errCloseRegisters != nil {
		log.Println("FixedExpenses: errCloseRegisters: " + errCloseRegisters.Error())
	}

	encoder := json.NewEncoder(w)
	encoder.Encode(fixedExpenses)
}

func addFixedExpense(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	var newFixedExpense FixedExpenses
	json.Unmarshal(body, &newFixedExpense)

	// insert in database
	result, errInsert := db2.Db.Exec("INSERT INTO fixed_expenses (name_expense,value_expense,due_date,pay_date) VALUES (?,?,?,?)", newFixedExpense.NameExpense, newFixedExpense.ValueExpense, newFixedExpense.DueDate, newFixedExpense.PayDate)

	idGenerated, errLastInsertId := result.LastInsertId()

	if errInsert != nil || errLastInsertId != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	newFixedExpense.Idfixed = int(idGenerated)

	w.WriteHeader(http.StatusCreated)
	encoder := json.NewEncoder(w)
	encoder.Encode(newFixedExpense)

}

func ConfigRoute(router *mux.Router) {
	// router.HandleFunc("/", test)
	router.HandleFunc("/fixedexpenses", listFixedExpenses).Methods("GET")
	router.HandleFunc("/fixedexpenses", addFixedExpense).Methods("POST")
}
