# Asaas CLI

A command-line tool for managing Asaas webhooks.

## Installation

```bash
go build -o asaascli cmd/main.go
```

## Configuration

The CLI uses a configuration file stored in your home directory (`~/.asaascli.conf`). You need to set up your API key before using the tool.

### Setting up your API key

```bash
./asaascli config set-token <your-api-key>
```

### Setting environment (sandbox vs production)

By default, the CLI uses the production environment. To switch to sandbox:

```bash
./asaascli config set-sandbox true
```

To switch back to production:

```bash
./asaascli config set-sandbox false
```

### Viewing current configuration

```bash
./asaascli config show
```

This will display:
- Config file location
- API key (masked for security)
- Current environment setting

## Commands

### List webhooks

```bash
./asaascli list
```

### Update webhook URL

```bash
./asaascli update-webhook-url <webhook-id> --newurl <new-url>
```

### Toggle webhook sync

```bash
./asaascli toggle-webhook-sync <webhook-id> <true|false>
```

## Configuration File

The configuration file is stored as JSON in your home directory:

```json
{
  "api_key": "your-api-key-here",
  "sandbox": false
}
```

## Error Handling

- If no API key is set, the CLI will show an error message with instructions on how to set it
- If the configuration file is corrupted, the CLI will create a new one with default values
- All API errors are displayed with descriptive messages

## Security

- The configuration file is created with restricted permissions (600)
- API keys are masked when displayed in the configuration
- The configuration file is stored in your home directory for security 