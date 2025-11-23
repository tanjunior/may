package main

import (
	"context"
	"fmt"

	"gorm.io/gorm"
)

var database = db()
var ctx = context.Background()

func getLatestProduct() (Product, error) {

	product, err := gorm.G[Product](database).First(ctx) // find product with integer primary key
	if err != nil {
		fmt.Println("Error fetching product:", err)
		return Product{}, err
	}
	return product, err
}
