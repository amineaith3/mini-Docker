package main

import (
	"fmt"
	"os"

	"mini-docker/internal/manager"
)

func main() {
	manager.Run(os.Args)
}

func handle(err error) {
	if err != nil {
		fmt.Println("error : ", err)
		os.Exit(1)
	}
}
