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

type VariableExpenses struct {
	Idvariable   int     `json:"id"`
	NameExpense  string  `json:"name"`
	ValueExpense float64 `json:"value"`
	DueDate      string  `json:"duedate"`
	PayDate      string  `json:"paydate"`
	DateRegister string  `json:"dateregister"`
}

func listVariableExpenses(w http.ResponseWriter, r *http.Request) {
	registers, errSelect := db2.Db.Query("SELECT * FROM variable_expenses")

	if errSelect != nil {
		log.Println("variable_expenses: " + errSelect.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	var variableExpenses []VariableExpenses = make([]VariableExpenses, 0)
	for registers.Next() {
		var variableExpense VariableExpenses
		errScan := registers.Scan(&variableExpense.Idvariable, &variableExpense.NameExpense, &variableExpense.ValueExpense, &variableExpense.DueDate, &variableExpense.PayDate, &variableExpense.DateRegister)
		if errScan != nil {
			log.Println("FixedExepenses: errScan: " + errScan.Error())
			continue
		}

		variableExpenses = append(variableExpenses, variableExpense)
	}

	errCloseRegisters := registers.Close()

	if errCloseRegisters != nil {
		log.Println("VariableExpenses: errCloseRegisters: " + errCloseRegisters.Error())
	}

	encoder := json.NewEncoder(w)
	encoder.Encode(variableExpenses)
}

func validateVariables(variableExpenses VariableExpenses) string {
	if len(variableExpenses.NameExpense) == 0 || len(variableExpenses.NameExpense) > 50 {
		return "O campo Autor precisa ter o mínimo de 1 caractere e máximo de 50 caracteres!"
	}

	// Não houve erro de validação
	return ""
}

func addVariableExpense(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	var newVariableExpense VariableExpenses
	json.Unmarshal(body, &newVariableExpense)

	// validate
	errValidate := validateVariables(newVariableExpense)
	if len(errValidate) > 0 {
		w.WriteHeader(http.StatusUnprocessableEntity)
		json.NewEncoder(w).Encode(ResponseErr{errValidate})
		return
	}

	// insert in database
	result, errInsert := db2.Db.Exec("INSERT INTO variable_expenses (name_expense,value_expense,due_date,pay_date, date_current) VALUES (?,?,?,?,?)", newVariableExpense.NameExpense, newVariableExpense.ValueExpense, newVariableExpense.DueDate, newVariableExpense.PayDate, newVariableExpense.DateRegister)

	idGenerated, errLastInsertId := result.LastInsertId()

	if errInsert != nil || errLastInsertId != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	newVariableExpense.Idvariable = int(idGenerated)

	w.WriteHeader(http.StatusCreated)
	encoder := json.NewEncoder(w)
	encoder.Encode(newVariableExpense)

}

func alterVariableExpense(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["idvariable"])

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
	}

	bodycap, errBody := ioutil.ReadAll(r.Body)

	if errBody != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	var alterVariable VariableExpenses
	errJson := json.Unmarshal(bodycap, &alterVariable)

	if errJson != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	register := db2.Db.QueryRow("SELECT idvariable, name_expense, value_expense, due_date, pay_date, date_register FROM variable_expenses WHERE idvariable = ?", id)
	var variableExpense VariableExpenses
	errScan := register.Scan(&variableExpense.Idvariable, &variableExpense.NameExpense, &variableExpense.ValueExpense, &variableExpense.DueDate, &variableExpense.PayDate, &variableExpense.DateRegister)

	if errScan != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	_, errExec := db2.Db.Exec("UPDATE variable_expenses SET name_expense = ?, value_expense = ?, due_date = ?, pay_date = ?, date_current = ? WHERE idfixed = ?", variableExpense.NameExpense, variableExpense.ValueExpense, variableExpense.DueDate, variableExpense.PayDate, variableExpense.DateRegister, id)

	if errExec != nil {
		log.Println("AlterVariableExpense: errExec: " + errExec.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(alterVariable)
}

func deleteVariable(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["idvariable"])

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	register := db2.Db.QueryRow("SELECT idvariable FROM variable_expenses WHERE idvariable = ?", id)
	var idOfVariable int
	errScan := register.Scan(&idOfVariable)

	if errScan != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	_, errExec := db2.Db.Exec("DELETE FROM variable_expenses WHERE idvariable = ?", id)

	if errExec != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func filterVariableDate(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, _ := strconv.Atoi(vars["dateregister"])

	registers, errSelect := db2.Db.Query("SELECT * FROM variable_expenses WHERE date_register = ?", id)

	if errSelect != nil {
		log.Println("variable_expenses: " + errSelect.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	defer registers.Close()
	var variableExpenses []VariableExpenses = make([]VariableExpenses, 0)

	for registers.Next() {
		var variableExpense VariableExpenses
		errScan := registers.Scan(&variableExpense.Idvariable, &variableExpense.NameExpense, &variableExpense.ValueExpense, &variableExpense.DueDate, &variableExpense.PayDate, &variableExpense.DateRegister)
		if errScan != nil {
			log.Println("FixedExepenses: errScan: " + errScan.Error())
			continue
		}

		variableExpenses = append(variableExpenses, variableExpense)
	}

	if len(variableExpenses) == 0 {
		log.Println("Fixed expenses empty")
		w.WriteHeader(http.StatusNotFound)
		return
	}

	errCloseRegisters := registers.Close()

	if errCloseRegisters != nil {
		log.Println("filter date: errCloseRegisters: " + errCloseRegisters.Error())
	}

	encoder := json.NewEncoder(w)
	encoder.Encode(variableExpenses)
}
