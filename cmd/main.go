package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/spf13/cobra"

	"github.com/gabrielforster/asaas-cli/internal/asaas"
)

var prodBaseUrl = "https://api.asaas.com/v3"
var devBaseUrl = "https://api-sandbox.asaas.com/v3"
var baseUrl = prodBaseUrl

var newUrl string
var sandbox bool

func main() {
	httpClient := &http.Client{}

	apikey := os.Getenv("ASAAS_API_KEY")

	var rootCmd = &cobra.Command{
		Use:   "asaas-cli <command>",
		Short: "A simple Go application with a command and an optional flag",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			command := args[0]

			fmt.Println("Executing command:", command, " with new value:", newUrl)
		},
	}

	rootCmd.PersistentFlags().BoolVar(&sandbox, "sandbox", false, "Use the sandbox environment")

	getClient := func() *asaas.AsaasWebhookClient {
		if sandbox {
			fmt.Println("Using sandbox environment")
			baseUrl = devBaseUrl
		}

		clientConfig := asaas.ClientConfig{
			BaseURL:    baseUrl,
			HTTPClient: httpClient,
			APIKey:     apikey,
		}

		return asaas.NewAsaasWebhookClient(clientConfig)
	}

	listCommand := &cobra.Command{
		Use:   "list",
		Short: "List all webhooks",
		Run: func(cmd *cobra.Command, args []string) {
			assasClient := getClient()
			webhooks, err := assasClient.ListWebhooks()
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error listing webhooks: %v\n", err)
				os.Exit(1)
			}

			for _, webhook := range webhooks {
				fmt.Printf("ID: %s, Name: %s\n", webhook.ID, webhook.Name)
			}
		},
	}

	updateWebhookUrlCommand := &cobra.Command{
		Use:   "update-webhook-url <webhook_id>",
		Short: "Update the URL of a webhook",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			assasClient := getClient()
			webhookID := args[0]
			webhook, err := assasClient.UpdateWebhookURL(webhookID, newUrl)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error updating webhook URL: %v\n", err)
				os.Exit(1)
			}
			fmt.Printf("Webhook ID: %s updated with new URL: %s\n", webhook.ID, webhook.URL)
		},
	}
	updateWebhookUrlCommand.PersistentFlags().StringVar(&newUrl, "newurl", "default_value_if_not_provided", "new value for the command")

	toggleWebhookSyncCommand := &cobra.Command{
		Use:   "toggle-webhook-sync <webhook_id> <enabled>",
		Short: "Enable or disable a webhook sync queue",
		Args:  cobra.ExactArgs(2),
		Run: func(cmd *cobra.Command, args []string) {
			assasClient := getClient()
			webhookID := args[0]
			enabled := args[1] == "true"
			webhook, err := assasClient.ToggleWebhookSync(webhookID, enabled)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error toggling webhook: %v\n", err)
				os.Exit(1)
			}
			status := "disabled"
			if webhook.Enabled {
				status = "enabled"
			}

			fmt.Printf("Webhook ID: %s is now %s\n", webhook.ID, status)
		},
	}

	rootCmd.AddCommand(listCommand)
	rootCmd.AddCommand(updateWebhookUrlCommand)
	rootCmd.AddCommand(toggleWebhookSyncCommand)

	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
