package main

import (
	"fmt"
	"encoding/json"

)

type Product struct {
	ProductId string	`json:"productId,omitempty"`
	Name      string	`json:"Name,omitempty"`
	Quantity  int		`json:"Quantity,string,omitempty"`
	Brand     string	`json:"Brand,omitempty"`
}


func main(){

	const data = `{
		"productId": "P2",
		"Brand": "Toyota",
		"Name": "car",
		"Quantity": "5"
	}`

	// fmt.Println(data)
	
	product := Product{}
	
	json.Unmarshal([]byte(data), &product)
	
	fmt.Println("Add Product:")
	fmt.Println("ID:  ", product.ProductId)
	fmt.Println("Name: ", product.Name)
	fmt.Println("Brand:  ", product.Brand)
	fmt.Println("Quantity:", product.Quantity)
}