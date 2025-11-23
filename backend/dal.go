package main

import (
	"context"
	"fmt"

	"gorm.io/gorm"
)

var database *gorm.DB
var ctx = context.Background()

func getLatestProduct() (Product, error) {

	product, err := gorm.G[Product](database).First(ctx) // find product with integer primary key
	if err != nil {
		fmt.Println("Error fetching product:", err)
		return Product{}, err
	}
	return product, err
}

// getAllProducts returns a page of products and the total count.
func getAllProducts(page int, perPage int) ([]Product, int64, error) {
	var products []Product
	var total int64

	// Count total products
	if err := database.Model(&Product{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Apply pagination
	offset := 0
	if page > 0 {
		offset = (page - 1) * perPage
	}

	if err := database.Limit(perPage).Offset(offset).Find(&products).Error; err != nil {
		return nil, 0, err
	}

	return products, total, nil
}

func getProductByID(id string) (Product, error) {
	product, err := gorm.G[Product](database).Where("id = ?", id).Take(ctx) // find product with integer primary key
	if err != nil {
		fmt.Println("Error fetching product:", err)
		return Product{}, err
	}
	return product, err
}

// addProduct creates a product and returns it.
func addProduct(code string, price uint) (Product, error) {
	product := Product{Code: code, Price: price}
	result := database.Create(&product)
	return product, result.Error
}

// updateProduct updates fields of a product and returns the updated product.
func updateProduct(id uint, newCode string, newPrice uint) (Product, error) {
	var product Product
	result := database.First(&product, id)
	if result.Error != nil {
		return Product{}, result.Error
	}
	product.Code = newCode
	product.Price = newPrice
	result = database.Save(&product)
	return product, result.Error
}

func deleteProduct(id uint) error {
	result := database.Delete(&Product{}, id)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}
