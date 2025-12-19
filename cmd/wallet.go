/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"YSNP/internal"
	"fmt"
	"log"
	"os"

	"github.com/spf13/cobra"
)

var (
	private string
)

// walletCmd represents the wallet command
var walletCmd = &cobra.Command{
	Use:   "wallet",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("wallet called")
		filename := "output_encoder.jsonl"
		encBytes, err := os.ReadFile("output_encoder.json.enc")
		if err != nil {
			log.Fatal(err)
		}

		decryptedBytes, err := internal.DecryptJSONWithPrivateKey(private, encBytes)
		if err != nil {
			log.Fatal(err)
		}

		if err := os.WriteFile(filename, decryptedBytes, 0600); err != nil {
			log.Fatal(err)
		}
	},
}

func init() {
	walletCmd.Flags().StringVar(&private, "key", "k", "Wallet Key")
	rootCmd.AddCommand(walletCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// walletCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// walletCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
