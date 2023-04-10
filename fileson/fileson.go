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

	"github.com/jokkebk/fileson-go/util"
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

			// Get the file's modification time in UTC and size
			modTime := info.ModTime().UTC().Format("2006-01-02 15:04:05")
			size := info.Size()

			// Check if the file is in the map
			if entry, ok := fson[relPath]; ok {
				// Get the file's modification time and size from the map
				ftime := entry.(map[string]interface{})["modified_gmt"].(string)
				fsize := entry.(map[string]interface{})["size"].(float64)

				// Compare the file's modification time and size to the map's
				if modTime != ftime || size != int64(fsize) {
					fmt.Println("File has changed", relPath)
					fmt.Println(modTime, "vs.", ftime)
					fmt.Println(size, "vs.", fsize)

					// Recalculate the sha1 hash
					hash, err := util.CalculateSHA1(path)

					if err != nil {
						fmt.Printf("Could not calculate hash for %s: %s", relPath, err)
						os.Exit(1)
					}

					// Update the map
					fson[relPath] = map[string]interface{}{
						"modified_gmt": modTime,
						"size":         size,
						"sha1":         hash,
					}

					// Print the new hash
					fmt.Println(relPath, "new hash", hash, "vs.", entry.(map[string]interface{})["sha1"].(string))
				}
			} else { // If the file is not in the map, print it
				fmt.Println(relPath, "not found")
			}
		}
		return nil
	})
}
