/*
Package cmd

Copyright © 2025 Alexandre Alfa <linkedin.com/in/alexandrealfa>
*/
package cmd

import (
	"YSNP/pkg/entity"
	"bufio"
	"encoding/json"
	"fmt"
	"github.com/fatih/color"
	"github.com/spf13/cobra"
	"log"
	"os"
	"strings"
	"time"
)

var (
	length    int
	symbols   bool
	store     bool
	vaultName string
)

type PassSchema struct {
	Name string `json:"name"`
	Pass string `json:"pass"`
	Time string `json:"time"`
}

// passwordCmd represents the password command
var passwordCmd = &cobra.Command{
	Use:   "pa",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		password := entity.NewPassword(length, symbols)
		success := color.New(color.Bold, color.FgGreen).SprintFunc()
		fmt.Printf("Generated password:  %s\n", success(password))

		if store {
			reader := bufio.NewReader(os.Stdin)

			fmt.Print("Enter a name for this password: ")

			name, err := reader.ReadString('\n')
			if err != nil {
				log.Fatal("Error reading name:", err)
			}

			name = strings.TrimSpace(name)

			if name == "" {
				log.Fatal("Password name cannot be empty")
			}

			schema, err := newPassword(name, password)
			if err != nil {
				log.Fatal("error to create Password Schema: ", err)
			}

			filename := "output_encoder.jsonl"

			if vaultName != "" {
				filename = vaultName
			}

			schema.save(filename)

			fmt.Println("Password saved in vault.")
		}
	},
}

func newPassword(name, pass string) (*PassSchema, error) {
	return &PassSchema{
		name,
		pass,
		time.Now().Format(time.DateTime),
	}, nil
}

func (p *PassSchema) save(vaultName string) {
	p.Time = time.Now().Format(time.DateTime)

	file, err := os.OpenFile(
		vaultName,
		os.O_CREATE|os.O_WRONLY|os.O_APPEND,
		0644,
	)
	if err != nil {
		log.Fatal("error to open file:", err)
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")

	if err := encoder.Encode(p); err != nil {
		log.Fatal("Error encoding JSON to file:", err)
	}
}

func init() {
	passwordCmd.Flags().IntVarP(&length, "length", "l", 10, "Password length")
	passwordCmd.Flags().BoolVarP(&symbols, "symbols", "s", false, "Include symbols")
	passwordCmd.Flags().BoolVar(&store, "store", false, "Store password in vault")

	rootCmd.AddCommand(passwordCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// passwordCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// passwordCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
