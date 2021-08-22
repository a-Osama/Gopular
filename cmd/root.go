package cmd

import (
	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "gopular",
	Short: "A CLI tool to find public repos",
	Long: `Get the most starred public repositories of your favorite programming language:

Gopular helps you find the most popular repositories based on criteria like
the programming language of choice, created after a specific date, and how many you want.`,
}

func Execute() {
	cobra.CheckErr(rootCmd.Execute())
}
