package main

import (
	"errors"
	"fmt"
	"github.com/manifoldco/promptui"
	"math"
	"strconv"
)

func main() {
	pretty()
	preferred := []string{"EUR", "USD", "MAD"}
	url := apiUrl()
	resp, err := ApiResponse(url)
	if err != nil {
		return
	}
	items := make([]string, 0, len(resp.Rates))
	for _, code := range preferred {
		if _, ok := resp.Rates[code]; ok {
			items = append(items, code)
		}
	}

	for code := range resp.Rates {
		skip := false
		for _, p := range items {
			if code == p {
				skip = true
				break
			}
		}
		if !skip {
			items = append(items, code)
		}
	}

	fromCode := mFirst(items)
	toCode := mSecond(items)

	from, ok1 := GetRate(resp, fromCode)

	to, ok2 := GetRate(resp, toCode)

	if !ok1 || !ok2 {
		fmt.Println("somthing went wrong")
		return
	}
	amount, err := getAmount()
	if err != nil {
		return
	}
	result := (amount / from) * to
	fmt.Printf("%.2f %s = %.2f %s", amount, fromCode, result, toCode)
}
func getAmount() (float64, error) {
	validate := func(input string) error {
		value, err := strconv.ParseFloat(input, 64)
		if err != nil {
			return errors.New("please enter a valid number")
		}
		if math.IsNaN(value) || math.IsInf(value, 0) {
			return errors.New("Invalid number")
		}
		if value < 0 {
			return errors.New("Number must be positive")
		}
		return nil
	}
	prompt := promptui.Prompt{
		Label:    "Amount",
		Validate: validate,
	}
	input, err := prompt.Run()
	if err != nil {
		return 0, fmt.Errorf("prompt failed: %w", err)

	}
	amount, _ := strconv.ParseFloat(input, 64)
	return amount, nil
}
func prompt(label string, items []string) string {
	prompt := promptui.Select{
		Label: label,
		Items: items,
	}

	_, r, err := prompt.Run()

	if err != nil {
		fmt.Printf("Prompt failed %v\n", err)
		return ""
	}
	return r
}
func mFirst(items []string) string {
	label := "enter your first currency "
	return prompt(label, items)

}
func mSecond(items []string) string {
	label := "enter your second currency "
	return prompt(label, items)

}
