package main

import "fmt"

type Calculation struct {
	Id         string `json:"id"`
	Expression string `json:"expretion"`
	Result     string `json:"result"`
}

type CalculationRequest struct {
	Expression string `json:"expression"`
}

var calculations = []Calculation{}

func main() {
	fmt.Println("hello")
}
