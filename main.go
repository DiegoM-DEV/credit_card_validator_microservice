package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

type Card struct {
	CardNumber string
}

type Response struct {
	IsValid bool `json:"verify"`
}

func LuhnAlgorithm(cardNumber string) bool {
	const dim int = 16

	var sum int = 0
	var isValid bool = false

	var count int = 0
	var cardNumberReverse [dim]int
	for i := len(cardNumber) - 1; i >= 0; i-- {
		var firtDigit int
		var secondDigit int
		digit, err := strconv.Atoi(cardNumber[count : count+1])
		count++
		if err == nil {
			if (i+1)%2 == 0 {
				var tmp int = digit * 2
				if tmp >= 10 {
					slice := strconv.Itoa(tmp)
					first, err := strconv.Atoi(slice[0:1])
					if err == nil {
						firtDigit = first
					}
					second, err := strconv.Atoi(slice[1:2])
					if err == nil {
						secondDigit = second
					}
					cardNumberReverse[i] = firtDigit + secondDigit
					sum = sum + cardNumberReverse[i]
				} else if tmp < 10 {
					cardNumberReverse[i] = tmp
					sum = sum + cardNumberReverse[i]
				}
			} else {
				cardNumberReverse[i] = digit
				sum = sum + cardNumberReverse[i]
			}
		}
	}

	if sum%10 == 0 {
		isValid = true
	} else if sum%10 != 0 {
		isValid = false
	}
	fmt.Println(cardNumberReverse)
	fmt.Println(sum)
	fmt.Println(isValid)
	return isValid
}

func verifyCard(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Endpoint verifyCard")
	var card Card
	err := json.NewDecoder(r.Body).Decode(&card)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	fmt.Println(card.CardNumber)
	isValid := LuhnAlgorithm(card.CardNumber)
	var res Response
	res.IsValid = isValid
	result, err := json.Marshal(res)
	if err != nil {
		log.Fatal(err)
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(result)
}

func handleRequest() {
	r := mux.NewRouter()
	r.HandleFunc("/card/verify", verifyCard).Methods("POST").Schemes("http")
	http.ListenAndServe(":3002", r)
}

func main() {
	handleRequest()
}
