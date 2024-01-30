/*
Copyright © 2024 zvirgilx
*/
package cmd

import (
	"context"
	"encoding/json"

	"github.com/spf13/cobra"
	"github.com/zvirgilx/searxng-go/kernel/internal/engine"
	"github.com/zvirgilx/searxng-go/kernel/internal/search"
)

// searchCmd represents the search command
var searchCmd = &cobra.Command{
	Use:   "search",
	Short: "Perform search operations on command line",
	Run: func(cmd *cobra.Command, args []string) {
		runSearch(cmd, args)
	},
	Args: cobra.ExactArgs(1),
}

func init() {
	rootCmd.AddCommand(searchCmd)
}

func runSearch(cmd *cobra.Command, args []string) {
	q := args[0]
	r := search.Search(context.Background(), engine.Options{Query: q})
	d, _ := json.MarshalIndent(r, "", "  ")
	cmd.Println(string(d))
}
