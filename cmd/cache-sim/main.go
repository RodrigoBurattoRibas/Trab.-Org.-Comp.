package main

import (
	"fmt"
	"os"

	"cache-simulator/internal/cache"
)

func main() {
	if err := cache.RunInteractiveMenu(os.Stdin, os.Stdout); err != nil {
		fmt.Fprintln(os.Stderr, "Erro:", err)
		os.Exit(1)
	}
}
