package model

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"
)

const (
	host     = "localhost"
	port     = 5432
	user     = "postgres"
	password = "123"
	database = "productsdb"
)

var db *sql.DB

func init() {
	var err error
	connString := fmt.Sprintf("host=%s port=%d user=%s password=%s database=%s sslmode=disable", host, port, user, password, database)

	db, err = sql.Open("postgres", connString)
	if err != nil {
		log.Fatal(err)
	}
}

type Product struct {
	ID                 int `json:"id"`
	Title, Description string
	Price              float32
}

func InsertProduct(data Product) {

	result, err := db.Exec("INSERT INTO products(title,description,price) VALUES($1,$2,$3)", data.Title, data.Description, data.Price)
	if err != nil {
		log.Fatal(err)
	}
	rowsAffected, _ := result.RowsAffected()
	fmt.Printf("Etkilenen kayıt sayısı(%d)", rowsAffected)
}

func UpdateProduct(data Product) {

	result, err := db.Exec("UPDATE product SET title=$2, description=$3, price =$4 WHERE id=$1)", data.ID, data.Title, data.Description, data.Price)
	if err != nil {
		log.Fatal(err)
	}
	rowsAffected, _ := result.RowsAffected()
	fmt.Printf("Etkilenen kayıt sayısı(%d)", rowsAffected)
}

func GetProducts() {
	rows, err := db.Query("SELECT * FROM products")
	if err != nil {
		if err == sql.ErrNoRows {
			fmt.Println("No Records Found!")
			return
		}
		log.Fatal(err)
	}
	defer rows.Close()

	var products []*Product
	for rows.Next() {
		prd := &Product{}
		err := rows.Scan(&prd.ID, &prd.Title, &prd.Description, &prd.Price)
		if err != nil {
			log.Fatal(err)
		}
		products = append(products, prd)
	}
	if err = rows.Err(); err != nil {
		log.Fatal(err)
	}
	for _, pr := range products {
		fmt.Printf("%d-%s,%s,$%.2f\n,", pr.ID, pr.Description, pr.Price, pr.Title)
	}
}

func GetProductByID(id int) {

	var product string
	err := db.QueryRow("SELECT title FROM products WHERE id=$1", id).Scan(&product)
	switch {
	case err == sql.ErrNoRows:
		log.Println("No product with that ID.")
	case err != nil:
		log.Fatal(err)
	default:
		fmt.Printf("Product is %s\n", product)
	}

}
