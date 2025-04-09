package main

import (
	"os"

	"github.com/bohdan-vykhovanets/url-shortener-svc/internal/cli"
)

func main() {
	if !cli.Run(os.Args) {
		os.Exit(1)
	}
}
