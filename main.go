package main

import (
	"fmt"
	"os"
	"path/filepath"
)

func main() {
	// Check that a file name was provided
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run main.go <filename> <folder>")
		return
	}

	// Read the file
	fson, err := readFileson(os.Args[1])

	if err != nil {
		fmt.Println("Error reading file:", err)
		return
	}

	// Print the map length
	fmt.Println(len(fson), "objects read from", os.Args[1])

	// Get file info
	dirname := os.Args[2]
	// Throw error if it's not a directory
	stat, err := os.Stat(dirname)
	if err != nil || !stat.IsDir() {
		fmt.Println("Error reading directory:", err)
		return
	}

	// Walk the directory and print all files
	filepath.Walk(dirname, func(path string, info os.FileInfo, err error) error {
		if !info.IsDir() {
			// Extract relative path from absolute path
			relPath, err := filepath.Rel(dirname, path)
			if err != nil {
				fmt.Println("Error getting relative path:", err)
				return err
			}

			// Check if the file is in the map
			if _, ok := fson[relPath]; !ok {
				fmt.Println(relPath, "not found")
			}
		}
		return nil
	})
}
