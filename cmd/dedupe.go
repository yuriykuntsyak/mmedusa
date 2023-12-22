/*
Copyright © 2023 Yuriy Kuntsyak
*/
package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// dedupeCmd represents the dedupe command
var dedupeCmd = &cobra.Command{
	Use:   "dedupe",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("dedupe called")
	},
}

func init() {
	rootCmd.AddCommand(dedupeCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// dedupeCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// dedupeCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
