package main

import (
	"ybilaly/model"
)

func main() {
	product := model.Product{
		ID:          1,
		Title:       "Golang Book",
		Description: "It's a good book",
		Price:       42.12,
	}
	model.InsertProduct(product)
}
