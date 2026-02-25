// Package manager is used to manage containers
package manager

import (
	"fmt"

	"mini-docker/internal/runtime"
)

func Run(args []string) {
	switch args[1] {
	case "run", "child":
		// generate the json
		runtime.Run(args)
	case "ps":
	case "inspect":
	default:
		fmt.Println("We will do something soon")
	}
}
