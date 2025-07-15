package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/spf13/cobra"

	"github.com/gabrielforster/asaas-cli/internal/asaas"
)

var prodBaseUrl = "https://www.asaas.com/api/v3"
var devBaseUrl = "https://sandbox.asaas.com/api/v3"
var baseUrl = prodBaseUrl

var newValue string
var sandbox bool

func main() {
	httpClient := &http.Client{}

	apikey := os.Getenv("ASAAS_API_KEY")

	var rootCmd = &cobra.Command{
		Use:   "myapp <command>",
		Short: "A simple Go application with a command and an optional flag",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			command := args[0]

			fmt.Println("Executing command:", command, " with new value:", newValue)
		},
	}

	rootCmd.PersistentFlags().StringVarP(&newValue, "newurl", "n", "default_value_if_not_provided", "new value for the command")
	rootCmd.PersistentFlags().BoolVar(&sandbox, "sandbox", false, "whether to use the sandbox environment")

	if sandbox {
		baseUrl = devBaseUrl
	}

	clientConfig := asaas.ClientConfig{
		BaseURL:    devBaseUrl,
		HTTPClient: httpClient,
		APIKey:     apikey,
	}

	assasClient := asaas.NewAsaasWebhookClient(clientConfig)

	rootCmd.AddCommand(&cobra.Command{
		Use:   "list",
		Short: "List all webhooks",
		Run: func(cmd *cobra.Command, args []string) {
			webhooks, err := assasClient.ListWebhooks()
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error listing webhooks: %v\n", err)
				os.Exit(1)
			}

			for _, webhook := range webhooks {
				fmt.Printf("ID: %s, Name: %s\n", webhook.ID, webhook.Name)
			}
		},
	})

	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
