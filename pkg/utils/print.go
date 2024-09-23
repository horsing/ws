// Package utils contains all kinds of utilities
package utils

import (
	"fmt"
	"os"
)

func isVerbose() bool {
	verbose := false
	for _, a := range os.Args {
		switch a {
		case "-V", "--verbose":
			verbose = true
		}
	}
	return verbose
}

var verbose = isVerbose()

// Print equivalent to fmt.Printf(format, args...)
func Print(format string, args ...any) {
	fmt.Fprintf(os.Stdout, format+"\n", args...)
}

// Print equivalent to fmt.Printf(format, args...)
func Log(format string, args ...any) {
	if verbose {
		fmt.Fprintf(os.Stdout, format+"\n", args...)
	}
}

// Error equivalent to fmt.Printf(format, args...)
func Error(err error) {
	fmt.Println(err)
}
