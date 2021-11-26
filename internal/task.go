package internal

import (
	"fmt"
	"os"
)

func PrintTask() error {
	r, err := os.ReadFile("README.md")

	if err != nil {
		return err
	}

	fmt.Println()
	fmt.Println(string(r))

	return nil
}
