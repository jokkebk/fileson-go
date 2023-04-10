/*
Copyright © Joonas Pihlajamaa <joonas.pihlajamaa@iki.fi>
*/
package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "fileson-go",
	Short: "Fileson backup tool with JSON file database and logs",
	Long: `Fileson is a backup tool that uses a JSON file as a database. You can
scan directories and add files to the database, backup scanned files and restore them.`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	// Run: func(cmd *cobra.Command, args []string) { },
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().BoolP("verbose", "v", false, "Verbose output")
}
