/*
Copyright © Joonas Pihlajamaa <joonas.pihlajamaa@iki.fi>
*/
package cmd

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/jokkebk/fileson-go/util"
)

// sha1Cmd represents the sha1 command
var sha1Cmd = &cobra.Command{
	Use:   "sha1 [file]",
	Args:  cobra.ExactArgs(1),
	Short: "Calculate sha1 of a file",
	RunE: func(cmd *cobra.Command, args []string) error {
		// Assign argument to variable
		filename := args[0]
		hash, err := util.CalculateSHA1(filename)

		if err != nil {
			fmt.Println(err)
			return err
		}

		// Use same format as sha1sum
		fmt.Printf("%s  %s\n", hash, filename)

		return nil
	},
}

func init() {
	rootCmd.AddCommand(sha1Cmd)
}
