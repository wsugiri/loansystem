package main

import (
	"fmt"

	"github.com/wsugiri/loansystem/app"
)

func main() {
	fmt.Println("starting server")
	err := app.SetupAndRunApp()
	if err != nil {
		panic(err)
	}
}
