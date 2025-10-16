package main

import (
	"fmt"
	"slices"
)

func CurrencyList() []string {
	url := apiURL()
	resp, err := ApiResponse(url)
	if err != nil {
		return nil
	}
	preferred := []string{"EUR", "USD", "MAD"}
	items := make([]string, 0, len(resp.Rates))
	for _, code := range preferred {
		if _, ok := resp.Rates[code]; ok {
			items = append(items, code)
		}
	}

	for code := range resp.Rates {
		if !slices.Contains(items, code) {
			items = append(items, code)
		}
	}
	return items
}
func app() string {

	url := apiURL()
	resp, err := ApiResponse(url)
	if err != nil {
		return ""
	}
	fromCode := mFirst(CurrencyList())
	toCode := mSecond(CurrencyList())

	from, ok1 := GetRate(resp, fromCode)

	to, ok2 := GetRate(resp, toCode)

	if !ok1 || !ok2 {
		fmt.Println("somthing went wrong")
		return ""
	}
	amount, err := getAmount()
	if err != nil {
		return ""
	}
	result := (amount / from) * to

	return fmt.Sprintf("%.2f %s = %.2f %s", amount, fromCode, result, toCode)

}
