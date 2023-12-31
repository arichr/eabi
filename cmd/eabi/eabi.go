package main

import (
	"fmt"
	"os"

	"github.com/arichr/eabi/pkg/eabi"
)

func main() {
	contents := []any{2, int16(255), int32(-257), uint16(259), nil}

	data, err := eabi.Marshal(contents)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when marshaling data: %s\n", err)
		os.Exit(1)
	}

	file, err := os.Create("out.private.eabi")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when opening a file: %s\n", err)
		os.Exit(1)
	}
	if _, err := file.Write(data); err != nil {
		fmt.Fprintf(os.Stderr, "Error when writing to the file: %s\n", err)
		os.Exit(1)
	}

	println("OK")
}
