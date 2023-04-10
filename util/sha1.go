/*
Copyright © Joonas Pihlajamaa <joonas.pihlajamaa@iki.fi>

This file contains simple utility to calculate sha1 sum of a file.
*/
package util

import (
	"crypto/sha1"
	"encoding/hex"
	"io"
	"os"
)

func CalculateSHA1(filename string) (string, error) {
	// Open file for reading, name from command line
	file, err := os.Open(filename)
	if err != nil {
		return "", err
	}
	defer file.Close()

	// Create new sha1 hash
	hash := sha1.New()

	// Copy file contents to hash
	_, err = io.Copy(hash, file)
	if err != nil {
		return "", err
	}

	// Get the 20 bytes hash
	hashInBytes := hash.Sum(nil)[:20]

	// Convert the bytes to a string
	return hex.EncodeToString(hashInBytes), nil
}
