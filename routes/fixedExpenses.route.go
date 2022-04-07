package routes

import (
	"context"
	"encoding/json"
	db2 "financasApi/db"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

var ctx context.Context

type FixedExpenses struct {
	Idfixed      int     `json:"id"`
	NameExpense  string  `json:"name"`
	ValueExpense float64 `json:"value"`
	DueDate      string  `json:"duedate"`
	PayDate      string  `json:"paydate"`
	DateCurrent  string  `json:"datecurrent"`
}

type ResponseErr struct {
	Erro string `json:"erro"`
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
		errScan := registers.Scan(&fixedExpense.Idfixed, &fixedExpense.NameExpense, &fixedExpense.ValueExpense, &fixedExpense.DueDate, &fixedExpense.PayDate, &fixedExpense.DateCurrent)
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

func validate(fixedExpenses FixedExpenses) string {
	if len(fixedExpenses.NameExpense) == 0 || len(fixedExpenses.NameExpense) > 50 {
		return "O campo Autor precisa ter o mínimo de 1 caractere e máximo de 50 caracteres!"
	}

	// Não houve erro de validação
	return ""
}

func addFixedExpense(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	var newFixedExpense FixedExpenses
	json.Unmarshal(body, &newFixedExpense)

	// validate
	errValidate := validate(newFixedExpense)
	if len(errValidate) > 0 {
		w.WriteHeader(http.StatusUnprocessableEntity)
		json.NewEncoder(w).Encode(ResponseErr{errValidate})
		return
	}

	// insert in database
	result, errInsert := db2.Db.Exec("INSERT INTO fixed_expenses (name_expense,value_expense,due_date,pay_date, date_current) VALUES (?,?,?,?,?)", newFixedExpense.NameExpense, newFixedExpense.ValueExpense, newFixedExpense.DueDate, newFixedExpense.PayDate, newFixedExpense.DateCurrent)

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

func alterFixedExpense(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["idfixed"])

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
	}

	bodycap, errBody := ioutil.ReadAll(r.Body)

	if errBody != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	var alterFixed FixedExpenses
	errJson := json.Unmarshal(bodycap, &alterFixed)

	if errJson != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	register := db2.Db.QueryRow("SELECT idfixed, name_expense, value_expense, due_date, pay_date FROM fixed_expenses WHERE idfixed = ?", id)
	var fixedExpense FixedExpenses
	errScan := register.Scan(&fixedExpense.Idfixed, &fixedExpense.NameExpense, &fixedExpense.ValueExpense, &fixedExpense.DueDate, &fixedExpense.PayDate, &fixedExpense.DateCurrent)

	if errScan != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	_, errExec := db2.Db.Exec("UPDATE fixed_expenses SET name_expense = ?, value_expense = ?, due_date = ?, pay_date = ?, date_current = ? WHERE idfixed = ?", alterFixed.NameExpense, alterFixed.ValueExpense, alterFixed.DueDate, alterFixed.PayDate, alterFixed.DateCurrent, id)

	if errExec != nil {
		log.Println("AlterFixedExpense: errExec: " + errExec.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(alterFixed)
}

func deleteFixed(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["idfixed"])

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	register := db2.Db.QueryRow("SELECT idfixed FROM fixed_expenses WHERE idfixed = ?", id)
	var idOfExpense int
	errScan := register.Scan(&idOfExpense)

	if errScan != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	_, errExec := db2.Db.Exec("DELETE FROM fixed_expenses WHERE idfixed = ?", id)

	if errExec != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func filterDate(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, _ := strconv.Atoi(vars["datecurrent"])
	// id := vars["datecurrent"]

	registers, errSelect := db2.Db.Query("SELECT * FROM fixed_expenses WHERE date_current = ?", id)

	if errSelect != nil {
		log.Println("fixed_expenses: " + errSelect.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	defer registers.Close()
	var fixedExpenses []FixedExpenses = make([]FixedExpenses, 0)

	for registers.Next() {
		var fixedExpense FixedExpenses
		errScan := registers.Scan(&fixedExpense.Idfixed, &fixedExpense.NameExpense, &fixedExpense.ValueExpense, &fixedExpense.DueDate, &fixedExpense.PayDate, &fixedExpense.DateCurrent)
		if errScan != nil {
			log.Println("FixedExepenses: errScan: " + errScan.Error())
			continue
		}

		fixedExpenses = append(fixedExpenses, fixedExpense)
	}

	if len(fixedExpenses) == 0 {
		log.Println("Fixed expenses empty")
		w.WriteHeader(http.StatusNotFound)
		return
	}

	errCloseRegisters := registers.Close()

	if errCloseRegisters != nil {
		log.Println("filter date: errCloseRegisters: " + errCloseRegisters.Error())
	}

	// errScan := register.Scan(&fixedExpense.Idfixed, &fixedExpense.NameExpense, &fixedExpense.ValueExpense, &fixedExpense.DueDate, &fixedExpense.PayDate, &fixedExpense.DateCurrent)

	// if errScan != nil {
	// 	log.Println("filterDate: errScan: " + errScan.Error())
	// 	w.WriteHeader(http.StatusNotFound)
	// 	return
	// }

	encoder := json.NewEncoder(w)
	encoder.Encode(fixedExpenses)
}

// func test(w http.ResponseWriter, r *http.Request) {
// 	fmt.Fprint(w, "teste")
// }

func ConfigRoute(router *mux.Router) {
	// router.HandleFunc("/", test)
	router.HandleFunc("/fixedexpenses", listFixedExpenses).Methods("GET")
	router.HandleFunc("/addfixedexpense", addFixedExpense).Methods("POST")
	router.HandleFunc("/fixedexpenses/{idfixed}", alterFixedExpense).Methods("PUT")
	router.HandleFunc("/fixedexpenses/{idfixed}", deleteFixed).Methods("DELETE")
	router.HandleFunc("/fixedexpensesdate/{datecurrent}", filterDate).Methods("GET")
}
