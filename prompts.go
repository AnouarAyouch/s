package main

import (
	"errors"
	"fmt"
	"math"
	"os"
	"strconv"

	"github.com/manifoldco/promptui"
)

func PropmtContinue() (string, error) {

	prompt := promptui.Prompt{
		Label: "Do you want to continue (yes/y)?",
	}
	input, err := prompt.Run()
	if err != nil {
		if err == promptui.ErrInterrupt {
			os.Exit(0)
		}
		fmt.Printf("Prompt failed %v\n", err)

	}

	return input, nil

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
		if err == promptui.ErrInterrupt {
			os.Exit(0)
		}
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

		if err == promptui.ErrInterrupt {
			os.Exit(0)
		}
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
