package database

import (
	"database/sql"
	"fmt"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/pianisimo/go-postgresql-pq/models"
	"log"
	"os"
)

func createConnection() *sql.DB {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	dataSourceName := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		os.Getenv("DB_HOST"), os.Getenv("DB_PORT"), os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"), os.Getenv("DB_NAME"), os.Getenv("DB_SSL"))

	db, err := sql.Open("postgres", dataSourceName)
	if err != nil {
		panic(err.(any))
	}

	err = db.Ping()
	if err != nil {
		panic(err.(any))
	}

	fmt.Println("Successfully connected to postgres")

	return db
}

func CreateStock(stock models.Stock) int64 {
	db := createConnection()
	defer db.Close()

	sqlStatement := `INSERT INTO stocks (name, price, company) VALUES ($1, $2, $3) RETURNING stockid`
	var id int64

	err := db.QueryRow(sqlStatement, stock.Name, stock.Price, stock.Company).Scan(&id)
	if err != nil {
		log.Fatalf("unable to execute query. %v", err)
	}

	fmt.Printf("inserted a single record %v", id)
	return id
}

func GetStock(id int64) (models.Stock, error) {
	db := createConnection()
	defer db.Close()

	var stock models.Stock

	sqlStatement := `SELECT * FROM stocks WHERE stockid=$1`

	row := db.QueryRow(sqlStatement, id)
	err := row.Scan(&stock.StockId, &stock.Name, &stock.Price, &stock.Company)

	switch err {
	case sql.ErrNoRows:
		fmt.Println("no rows were returned")
		return stock, nil
	case nil:
		return stock, nil
	default:
		log.Fatalf("unable to scan row. %v", err)
	}

	return stock, err
}

func GetAllStocks() ([]models.Stock, error) {
	db := createConnection()
	defer db.Close()

	var stocks []models.Stock

	sqlStatement := `SELECT * FROM stocks`
	rows, err := db.Query(sqlStatement)
	if err != nil {
		log.Fatalf("unable execute query. %v", err)
	}
	defer rows.Close()
	for rows.Next() {
		var stock models.Stock
		err = rows.Scan(&stock.StockId, &stock.Name, &stock.Price, &stock.Company)
		if err != nil {
			log.Fatalf("unable to scan row. %v", err)
		}

		stocks = append(stocks, stock)
	}

	return stocks, err
}

func UpdateStock(id int64, stock models.Stock) int64 {
	db := createConnection()
	defer db.Close()

	sqlStatement := `UPDATE stocks SET name=$2, price=$3, company=$4 WHERE stockid=$1`

	res, err := db.Exec(sqlStatement, id, stock.Name, stock.Price, stock.Company)
	if err != nil {
		log.Fatalf("unable to execute query, %v", err)
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		log.Fatalf("error while checking affected rows, %v", err)
	}

	fmt.Printf("totla row/records affected. %v", rowsAffected)
	return rowsAffected
}

func DeleteStock(id int64) int64 {
	db := createConnection()
	defer db.Close()

	sqlStatement := `DELETE FROM stocks WHERE stockid=$1`

	res, err := db.Exec(sqlStatement, id)
	if err != nil {
		log.Fatalf("unable to execute query, %v", err)
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		log.Fatalf("error while checking affected rows, %v", err)
	}

	fmt.Printf("totla row/records affected. %v", rowsAffected)
	return rowsAffected
}
