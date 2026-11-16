package main

import (
	"fmt"
	"os"

	"tas-app/internal/infrastructure/container"
)

func main() {
	terminalCLI, err := container.BuildCLI()
	if err != nil {
		fmt.Fprintf(os.Stderr, "не удалось создать CLI: %v\n", err)
		os.Exit(1)
	}

	terminalCLI.Run()
}
