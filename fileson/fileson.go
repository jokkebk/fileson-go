/*
Copyright © Joonas Pihlajamaa <joonas.pihlajamaa@iki.fi>

Fileson is a map-like structure that records changes in JSON file.
All changes are logged and appended as they happen.
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

type Fileson struct {
	data map[string]interface{}
}

func NewFileson() *Fileson {
	return &Fileson{
		data: make(map[string]interface{}),
	}
}

func (m *Fileson) Len() int {
	return len(m.data)
}

func (m *Fileson) Delete(key string) {
	fmt.Printf("Deleting %s\n", key)
	delete(m.data, key)
}

func (m *Fileson) Range(f func(key string, value interface{}) bool) {
	for k, v := range m.data {
		if !f(k, v) {
			break
		}
	}
}

func (m *Fileson) Load(key string) (value interface{}, ok bool) {
	value, ok = m.data[key]
	return
}

func (m *Fileson) Store(key string, value interface{}) {
	fmt.Printf("Storing %s = %s\n", key, value)

	// Create array with key and value as two elements
	objects := []interface{}{key, value}
	jsonStr, err := json.Marshal(objects)

	if err != nil {
		fmt.Println("Error marshalling JSON:", err)
		return
	}

	fmt.Printf("JSON: %s\n", jsonStr)

	m.data[key] = value
}

func ReadFile(filename string) (*Fileson, error) {
	fileson := &Fileson{
		data: make(map[string]interface{}),
	}

	// Open file for reading, name from command line
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

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
				// Add key and value to map bypassing Store
				fileson.data[key] = objects[1]
			}
		} else if len(objects) == 1 {
			// If there was only one item, delete the key, bypassing Delete
			if key, ok := objects[0].(string); ok {
				delete(fileson.data, key)
			}
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return fileson, nil
}

func (fson *Fileson) ScanDirectory(dirname string) {
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
			if entry, ok := fson.Load(relPath); ok {
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

					// Print the new hash
					fmt.Println(relPath, "new hash", hash, "vs.", entry.(map[string]interface{})["sha1"].(string))

					// Update the map
					fson.Store(relPath, map[string]interface{}{
						"modified_gmt": modTime,
						"size":         size,
						"sha1":         hash,
					})
				}
			} else { // If the file is not in the map, print it
				fmt.Println(relPath, "not found")
			}
		}
		return nil
	})
}
