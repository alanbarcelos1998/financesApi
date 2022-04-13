package routes

import (
	"encoding/json"
	db2 "financasApi/db"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
)

type Expenses struct {
	Idexpense    int     `json:"id"`
	Name         string  `json:"name"`
	Type         string  `json:"type"`
	Value        float64 `json:"value"`
	PayDate      string  `json:"paydate"`
	RegisterDate string  `json:"registerdate"`
}

func listExpenses(w http.ResponseWriter, r *http.Request) {
	registers, errSelect := db2.Db.Query("SELECT * FROM expenses")

	if errSelect != nil {
		log.Println("Expenses: " + errSelect.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	var expenses []Expenses = make([]Expenses, 0)
	for registers.Next() {
		var expense Expenses
		errScan := registers.Scan(&expense.Idexpense, &expense.Name, &expense.Type, &expense.Value, &expense.PayDate, &expense.RegisterDate)
		if errScan != nil {
			log.Println("expenses: errScan: " + errScan.Error())
			continue
		}

		expenses = append(expenses, expense)
	}

	errCloseRegisters := registers.Close()

	if errCloseRegisters != nil {
		log.Println("FixedExpenses: errCloseRegisters: " + errCloseRegisters.Error())
	}

	encoder := json.NewEncoder(w)
	encoder.Encode(expenses)
}

func validate(expenses Expenses) string {
	if len(expenses.Name) == 0 || len(expenses.Name) > 60 {
		return "O campo Autor precisa ter o mínimo de 1 caractere e máximo de 50 caracteres!"
	}

	// Não houve erro de validação
	return ""
}

func cadExpense(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	var newExpense Expenses
	json.Unmarshal(body, &newExpense)

	// validate
	errValidate := validate(newExpense)
	if len(errValidate) > 0 {
		w.WriteHeader(http.StatusUnprocessableEntity)
		json.NewEncoder(w).Encode(ResponseErr{errValidate})
		return
	}

	// Put the current date and set in RegisterDate
	currentDate := time.Now()
	formatDate := currentDate.String()
	shortDate := formatDate[0:10]
	newExpense.RegisterDate = shortDate

	// insert in database
	result, errInsert := db2.Db.Exec("INSERT INTO expenses (expense_name,expense_type,value,pay_date, register_date) VALUES (?,?,?,?,?)", newExpense.Name, newExpense.Type, newExpense.Value, newExpense.PayDate, newExpense.RegisterDate)

	idGenerated, errLastInsertId := result.LastInsertId()

	if errInsert != nil || errLastInsertId != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	newExpense.Idexpense = int(idGenerated)

	w.WriteHeader(http.StatusCreated)
	encoder := json.NewEncoder(w)
	encoder.Encode(newExpense)

}

func alterExpense(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
	}

	bodycap, errBody := ioutil.ReadAll(r.Body)

	if errBody != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	var alter Expenses
	errJson := json.Unmarshal(bodycap, &alter)

	if errJson != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	register := db2.Db.QueryRow("SELECT idexpense, expense_name, expense_type, value, pay_date, register_date FROM expenses WHERE idexpense = ?", id)
	var expense Expenses
	errScan := register.Scan(&expense.Idexpense, &expense.Name, &expense.Type, &expense.Value, &expense.PayDate, &expense.RegisterDate)

	if errScan != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	_, errExec := db2.Db.Exec("UPDATE expenses SET expense_name = ?, expense_type = ?, value = ?, pay_date = ?, register_date = ? WHERE idexpense = ?", alter.Name, alter.Type, alter.Value, alter.PayDate, alter.RegisterDate, id)

	if errExec != nil {
		log.Println("Alter Expense: errExec: " + errExec.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(alter)
}

func delete(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	register := db2.Db.QueryRow("SELECT idexpense FROM expenses WHERE idexpense = ?", id)
	var idOfExpense int
	errScan := register.Scan(&idOfExpense)

	if errScan != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	_, errExec := db2.Db.Exec("DELETE FROM expenses WHERE idexpense = ?", id)

	if errExec != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)

}

func filterDate(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, _ := strconv.Atoi(vars["date"])

	registers, errSelect := db2.Db.Query("SELECT * FROM expenses WHERE register_date = ?", id)

	if errSelect != nil {
		log.Println("Expenses: " + errSelect.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	defer registers.Close()
	var expenses []Expenses = make([]Expenses, 0)

	for registers.Next() {
		var expense Expenses
		errScan := registers.Scan(&expense.Idexpense, &expense.Name, &expense.Type, &expense.Value, &expense.PayDate, &expense.RegisterDate)
		if errScan != nil {
			log.Println("Exepenses: errScan: " + errScan.Error())
			continue
		}

		expenses = append(expenses, expense)
	}

	if len(expenses) == 0 {
		log.Println("Fixed expenses empty")
		w.WriteHeader(http.StatusNotFound)
		return
	}

	errCloseRegisters := registers.Close()

	if errCloseRegisters != nil {
		log.Println("filter date: errCloseRegisters: " + errCloseRegisters.Error())
	}

	encoder := json.NewEncoder(w)
	encoder.Encode(expenses)
}
