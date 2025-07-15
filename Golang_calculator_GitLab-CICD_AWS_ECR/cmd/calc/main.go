package main

import (
	"calculator/internal/calculator"
	"fmt"
)

func main() {
	var four, five float64 = 4, 5
	fmt.Printf("Add %v to %v is %v\n", four, five, calculator.Add(four, five))
	fmt.Printf("Substract %v from %v is %v\n", five, four, calculator.Subtract(four, five))
	fmt.Printf("Multiply %v by %v is %v\n", four, five, calculator.Multiply(four, five))
	if result, err := calculator.Divide(four, five); err == nil {
		fmt.Printf("Divide %v by %v is %v\n", four, five, result)
	} else {
		fmt.Println("Divide error: ", err)
	}
	if result, err := calculator.Sqrt(five); err == nil {
		fmt.Printf("Sqrt of %v is %v\n", five, result)
	} else {
		fmt.Println("Sqrt error: ", err)
	}
}
