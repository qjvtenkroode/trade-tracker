package main

import (
	"fmt"
	"html/template"
	"log"
	"math/rand"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

// TODO add persistency through database (Postgress?)

type Trade struct {
    // TODO add type (e.g. turbo, stock)
    // TODO add underlying asset (AEX, Shell)
	ID           int
	Name         string
	Transactions []Transaction
}

type Transaction struct {
    // TODO add datetime
	ID     int
	Price  float64
	Amount int
    Costs float64
}

func main() {
	r := mux.NewRouter()
	baseTmpl := template.Must(template.ParseFiles("templates/index.html", "templates/navigation.html"))
	tradeTmpl := template.Must(template.ParseFiles("templates/tradeForm.html", "templates/navigation.html"))
	transactionTmpl := template.Must(template.ParseFiles("templates/transactionForm.html", "templates/navigation.html"))
	var listOfTrades []Trade

	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		baseTmpl.Execute(w, listOfTrades)
	}).Methods("GET")

	r.HandleFunc("/addtrade", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			tradeTmpl.Execute(w, nil)
			fmt.Fprintf(w, "len=%d cap=%d %v\n", len(listOfTrades), cap(listOfTrades), listOfTrades)
			return
		}

		trade := Trade{
			ID:   rand.Intn(1000),
			Name: r.FormValue("name"),
		}

		listOfTrades = append(listOfTrades, trade)
		tradeTmpl.Execute(w, struct{ Success bool }{true})
	}).Methods("GET", "POST")

	r.HandleFunc("/addtransaction", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			transactionTmpl.Execute(w, nil)
			return
		}

		tradeid, err := strconv.Atoi(r.FormValue("tradeid"))
		if err != nil {
			log.Fatal(err)
		}
		price, err := strconv.ParseFloat(r.FormValue("price"), 64)
		if err != nil {
			log.Fatal(err)
		}
		amount, err := strconv.Atoi(r.FormValue("amount"))
		if err != nil {
			log.Fatal(err)
		}
		costs, err := strconv.ParseFloat(r.FormValue("costs"), 64)
		if err != nil {
			log.Fatal(err)
		}

		transaction := Transaction{
			ID:     rand.Intn(1000),
			Price:  price,
			Amount: amount,
            Costs: costs,
		}

		listOfTrades[tradeid].Transactions = append(listOfTrades[tradeid].Transactions, transaction)

		transactionTmpl.Execute(w, struct{ Success bool }{true})
	}).Methods("GET", "POST")

	http.ListenAndServe(":8000", r)
}
