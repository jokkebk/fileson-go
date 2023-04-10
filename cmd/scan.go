/*
Copyright © Joonas Pihlajamaa <joonas.pihlajamaa@iki.fi>
*/
package cmd

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/jokkebk/fileson-go/fileson"
)

// scanCmd represents the scan command
var scanCmd = &cobra.Command{
	Use:   "scan [JSON file] [directory]",
	Args:  cobra.ExactArgs(2),
	Short: "Scan a directory for files and append/write to a JSON file",
	Long: `Traverse a directory and its subdirectories and add all
files to a JSON file. If the file already exists, append changes.`,
	Run: run,
}

func run(cmd *cobra.Command, args []string) {
	// Assign arguments to variables
	jsonFile := args[0]
	directory := args[1]

	// Read the file
	fson, err := fileson.ReadFile(jsonFile)

	if err != nil {
		fmt.Println("Error reading file:", err)
		return
	}

	// Print the map length
	fmt.Println(fson.Len(), "objects read from", jsonFile)

	fson.ScanDirectory(directory)
}

func init() {
	rootCmd.AddCommand(scanCmd)
}
