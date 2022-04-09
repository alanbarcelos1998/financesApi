package routes

import (
	"encoding/json"
	db2 "financasApi/db"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

type FixedExpenses struct {
	Idfixed      int     `json:"id"`
	NameExpense  string  `json:"name"`
	ValueExpense float64 `json:"value"`
	DueDate      string  `json:"duedate"`
	PayDate      string  `json:"paydate"`
	DateRegister string  `json:"dateregister"`
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
		errScan := registers.Scan(&fixedExpense.Idfixed, &fixedExpense.NameExpense, &fixedExpense.ValueExpense, &fixedExpense.DueDate, &fixedExpense.PayDate, &fixedExpense.DateRegister)
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

func validateFixed(fixedExpenses FixedExpenses) string {
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
	errValidate := validateFixed(newFixedExpense)
	if len(errValidate) > 0 {
		w.WriteHeader(http.StatusUnprocessableEntity)
		json.NewEncoder(w).Encode(ResponseErr{errValidate})
		return
	}

	// insert in database
	result, errInsert := db2.Db.Exec("INSERT INTO fixed_expenses (name_expense,value_expense,due_date,pay_date, date_register) VALUES (?,?,?,?,?)", newFixedExpense.NameExpense, newFixedExpense.ValueExpense, newFixedExpense.DueDate, newFixedExpense.PayDate, newFixedExpense.DateRegister)

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
	errScan := register.Scan(&fixedExpense.Idfixed, &fixedExpense.NameExpense, &fixedExpense.ValueExpense, &fixedExpense.DueDate, &fixedExpense.PayDate, &fixedExpense.DateRegister)

	if errScan != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	_, errExec := db2.Db.Exec("UPDATE fixed_expenses SET name_expense = ?, value_expense = ?, due_date = ?, pay_date = ?, date_current = ? WHERE idfixed = ?", alterFixed.NameExpense, alterFixed.ValueExpense, alterFixed.DueDate, alterFixed.PayDate, alterFixed.DateRegister, id)

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

func filterFixedDate(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, _ := strconv.Atoi(vars["dateregister"])

	registers, errSelect := db2.Db.Query("SELECT * FROM fixed_expenses WHERE date_register = ?", id)

	if errSelect != nil {
		log.Println("fixed_expenses: " + errSelect.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	defer registers.Close()
	var fixedExpenses []FixedExpenses = make([]FixedExpenses, 0)

	for registers.Next() {
		var fixedExpense FixedExpenses
		errScan := registers.Scan(&fixedExpense.Idfixed, &fixedExpense.NameExpense, &fixedExpense.ValueExpense, &fixedExpense.DueDate, &fixedExpense.PayDate, &fixedExpense.DateRegister)
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

	encoder := json.NewEncoder(w)
	encoder.Encode(fixedExpenses)
}
