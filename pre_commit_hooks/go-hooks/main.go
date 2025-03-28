package main

import (
	"fmt"
	"github.com/spf13/cobra"
)

var rootCmd = cobra.Command{
	Use:   "go-hooks",
	Short: "pre-commit hooks",
	Run:   func(cmd *cobra.Command, args []string) {},
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
	}
}
