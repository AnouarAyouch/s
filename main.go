package main

import "fmt"

func main() {
	for {
		output := app()
		fmt.Printf("%s\n", output)
		fmt.Println("")
		c, err := PropmtContinue()
		if err != nil {
			return
		}
		if c != "yes" && c != "y" {
			break
		}
	}

}
