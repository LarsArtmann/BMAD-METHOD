package main

import (
	"fmt"
	"os"

	"github.com/LarsArtmann/BMAD-METHOD/cmd/generator/commands"
)

const (
	// Version of the template generator
	Version = "1.0.0"
	
	// AppName is the application name
	AppName = "template-health-endpoint"
)

func main() {
	if err := commands.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
