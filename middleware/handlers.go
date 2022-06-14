package middleware

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/pianisimo/go-postgresql-pq/database"
	"github.com/pianisimo/go-postgresql-pq/models"
	"log"
	"net/http"
	"strconv"
)

type response struct {
	Id      int64  `json:"id,omitempty"`
	Message string `json:"message,omitempty"`
}

func CreateStock(w http.ResponseWriter, r *http.Request) {
	var stock models.Stock

	err := json.NewDecoder(r.Body).Decode(&stock)
	if err != nil {
		log.Fatalf("Can not decode to stock. %v", err)
	}

	insertId := database.CreateStock(stock)

	res := response{
		Id:      insertId,
		Message: "stock created successfully",
	}

	json.NewEncoder(w).Encode(res)
}

func GetStock(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id, err := strconv.Atoi(params["id"])
	if err != nil {
		log.Fatalf("unable to convert string into int. %v", err)
	}

	stock, err := database.GetStock(int64(id))
	if err != nil {
		log.Fatalf("unable to get stock. %v", err)
	}

	json.NewEncoder(w).Encode(stock)
}

func GetAllStock(w http.ResponseWriter, r *http.Request) {
	stocks, err := database.GetAllStocks()
	if err != nil {
		log.Fatalf("unable to get all the stocks. %v", err)
	}

	json.NewEncoder(w).Encode(stocks)
}

func UpdateStock(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id, err := strconv.Atoi(params["id"])
	if err != nil {
		log.Fatalf("unable to convert string into int. %v", err)
	}

	var stock models.Stock
	err = json.NewDecoder(r.Body).Decode(&stock)

	if err != nil {
		log.Fatalf("Can not decode to stock. %v", err)
	}

	updatedRows := database.UpdateStock(int64(id), stock)
	msg := fmt.Sprintf("stock updated successfully. Total rows/records affected %v", updatedRows)

	res := response{
		Id:      int64(id),
		Message: msg,
	}

	json.NewEncoder(w).Encode(res)
}

func DeleteStock(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id, err := strconv.Atoi(params["id"])
	if err != nil {
		log.Fatalf("unable to convert string into int. %v", err)
	}

	deletedRows := database.DeleteStock(int64(id))

	msg := fmt.Sprintf("stock deleted successfully. Total rows/records affected %v", deletedRows)

	res := response{
		Id:      int64(id),
		Message: msg,
	}

	json.NewEncoder(w).Encode(res)
}
