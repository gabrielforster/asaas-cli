package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/spf13/cobra"

	"github.com/gabrielforster/asaascli/internal/asaas"
	"github.com/gabrielforster/asaascli/internal/config"
)

var prodBaseUrl = "https://api.asaas.com/v3"
var devBaseUrl = "https://api-sandbox.asaas.com/v3"

var newUrl string

func main() {
	httpClient := &http.Client{}

	var rootCmd = &cobra.Command{
		Use:   "asaascli",
		Short: "Asaas CLI tool for managing webhooks",
		Long: `Asaas CLI is a command-line tool for managing Asaas webhooks.

Available commands:
  config              Manage configuration settings
  list                List all webhooks
  update-webhook-url  Update the URL of a webhook
  toggle-webhook-sync Enable or disable a webhook sync queue

Use 'asaascli <command> --help' for more information about a command.`,
	}

	getClient := func() (*asaas.AsaasWebhookClient, error) {
		apikey, err := config.GetAPIKey()
		if err != nil {
			return nil, err
		}

		isSandbox, err := config.IsSandbox()
		if err != nil {
			return nil, err
		}

		baseUrl := prodBaseUrl
		if isSandbox {
			fmt.Println("Using sandbox environment")
			baseUrl = devBaseUrl
		}

		clientConfig := asaas.ClientConfig{
			BaseURL:    baseUrl,
			HTTPClient: httpClient,
			APIKey:     apikey,
		}

		return asaas.NewAsaasWebhookClient(clientConfig), nil
	}

	configCommand := &cobra.Command{
		Use:   "config",
		Short: "Manage configuration settings",
	}

	setTokenCommand := &cobra.Command{
		Use:   "set-token <api-key>",
		Short: "Set the API key for Asaas API calls",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			apiKey := args[0]
			if err := config.SetAPIKey(apiKey); err != nil {
				fmt.Fprintf(os.Stderr, "Error setting API key: %v\n", err)
				os.Exit(1)
			}
			fmt.Println("API key set successfully")
		},
	}

	setSandboxCommand := &cobra.Command{
		Use:   "set-sandbox <true|false>",
		Short: "Set whether to use sandbox environment",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			sandboxStr := args[0]
			var sandbox bool
			switch sandboxStr {
			case "true":
				sandbox = true
			case "false":
				sandbox = false
			default:
				fmt.Fprintf(os.Stderr, "Error: sandbox must be 'true' or 'false'\n")
				os.Exit(1)
			}

			if err := config.SetSandbox(sandbox); err != nil {
				fmt.Fprintf(os.Stderr, "Error setting sandbox mode: %v\n", err)
				os.Exit(1)
			}

			env := "production"
			if sandbox {
				env = "sandbox"
			}
			fmt.Printf("Environment set to %s\n", env)
		},
	}

	showConfigCommand := &cobra.Command{
		Use:   "show",
		Short: "Show current configuration",
		Run: func(cmd *cobra.Command, args []string) {
			cfg, err := config.LoadConfig()
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error loading config: %v\n", err)
				os.Exit(1)
			}

			configPath, err := config.GetConfigPath()
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error getting config path: %v\n", err)
				os.Exit(1)
			}

			fmt.Printf("Config file: %s\n", configPath)
			fmt.Printf("API Key: %s\n", maskAPIKey(cfg.APIKey))
			fmt.Printf("Sandbox: %v\n", cfg.Sandbox)
		},
	}

	configCommand.AddCommand(setTokenCommand)
	configCommand.AddCommand(setSandboxCommand)
	configCommand.AddCommand(showConfigCommand)

	listCommand := &cobra.Command{
		Use:   "list",
		Short: "List all webhooks",
		Run: func(cmd *cobra.Command, args []string) {
			assasClient, err := getClient()
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error: %v\n", err)
				os.Exit(1)
			}

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
			assasClient, err := getClient()
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error: %v\n", err)
				os.Exit(1)
			}

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
			assasClient, err := getClient()
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error: %v\n", err)
				os.Exit(1)
			}

			webhookID := args[0]
			enabled := args[1] == "true"
			webhook, err := assasClient.ToggleWebhookSync(webhookID, enabled)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error toggling webhook: %v\n", err)
				os.Exit(1)
			}

			status := "disabled"
			if enabled {
				status = "enabled"
			}

			fmt.Printf("Webhook ID: %s is now %s\n", webhook.ID, status)
		},
	}

	rootCmd.AddCommand(configCommand)
	rootCmd.AddCommand(listCommand)
	rootCmd.AddCommand(updateWebhookUrlCommand)
	rootCmd.AddCommand(toggleWebhookSyncCommand)

	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

func maskAPIKey(apiKey string) string {
	if apiKey == "" {
		return "not set"
	}
	if len(apiKey) <= 8 {
		return "***"
	}
	return apiKey[:4] + "..." + apiKey[len(apiKey)-4:]
}
