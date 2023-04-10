/*
Copyright © Joonas Pihlajamaa <joonas.pihlajamaa@iki.fi>
*/
package fileson

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

// Function that returns a map with string keys and interface{} values
func ReadFileson(filename string) (map[string]interface{}, error) {
	// Open file for reading, name from command line
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	// Create map with string keys and interface{} values
	// This is the type of the JSON object we will read
	// from the file
	fileson := make(map[string]interface{})

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		jsonStr := scanner.Text()

		var objects []interface{}
		err := json.Unmarshal([]byte(jsonStr), &objects)
		if err != nil {
			return nil, err
		}

		// If there were no items in the array, skip
		if len(objects) == 0 {
			continue
		}

		// If the first item is not a string, panic
		if _, ok := objects[0].(string); !ok {
			return nil, fmt.Errorf("First item in array is not a string")
		}

		// Check if there were two items in array,
		// in which case the first is the key and the
		// second is the value
		if len(objects) == 2 {
			// Check that the first item is a string
			if key, ok := objects[0].(string); ok {
				// Add key and value to map
				fileson[key] = objects[1]
			}
		} else if len(objects) == 1 {
			// If there was only one item, delete the key
			if key, ok := objects[0].(string); ok {
				delete(fileson, key)
			}
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return fileson, nil
}

func ScanDirectory(dirname string, fson map[string]interface{}) {
	/*
			// Read the file
			fson, err := readFileson(os.Args[1])
		fileson

			if err != nil {
				fmt.Println("Error reading file:", err)
				return
			}

			// Print the map length
			fmt.Println(len(fson), "objects read from", os.Args[1])
	*/

	// Throw error if dirname is not a directory
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
