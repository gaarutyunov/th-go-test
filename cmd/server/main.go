package main

import (
	"th-go-test/internal"
	"th-go-test/pkg/hello"
)

func main() {
	hello.Println("server")
	if err := internal.PrintTask(); err != nil {
		panic(err)
	}
}
