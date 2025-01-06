package main

//go:generate go run ./internal/generator/generate.go code providers.json internal/libdns/registry.go
import (
	"context"
	"fmt"
	"os"
	"strings"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/project0/external-dns-webhook-libdns/internal/libdns"
	"github.com/project0/external-dns-webhook-libdns/internal/webhook"
	"github.com/urfave/cli/v3"
)

const (
	ENV_PREFIX                 = "WEBHOOK_LIBDNS_"
	flag_log_level             = "log-level"
	flag_log_format            = "log-format"
	flag_provider_name         = "provider-name"
	flag_provider_config       = "provider-config"
	flag_provider_config_file  = "provider-config-file"
	flag_zones                 = "domain-filter"
	flag_webhook_listen        = "webhook-listen"
	flag_webhook_read_timeout  = "webhook-read-timeout"
	flag_webhook_write_timeout = "webhook-write-timeout"
	flag_metrics_listen        = "metrics-listen"
)

func flagEnv(name string) string {
	return ENV_PREFIX + strings.ReplaceAll(strings.ToUpper(name), "-", "_")
}

func providerList() []string {
	var providers []string
	for k := range libdns.Registry {
		providers = append(providers, k)
	}
	return providers
}

func main() {

	cmd := &cli.Command{
		Name:  "external-dns-webhook-libdns",
		Usage: "Webhook for external-dns using libdns providers. Supported provider names: " + strings.Join(providerList(), ", "),

		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    flag_log_level,
				Usage:   "The log level (trace, debug, info, warn, error, fatal, panic)",
				Value:   "info",
				Sources: cli.EnvVars(flagEnv(flag_log_level)),
				Validator: func(s string) error {
					if _, err := log.ParseLevel(s); err != nil {
						return err
					}
					return nil
				},
			},
			&cli.StringFlag{
				Name:    flag_log_format,
				Usage:   "The log format (text, json)",
				Value:   "text",
				Sources: cli.EnvVars(flagEnv(flag_log_format)),
				Validator: func(s string) error {
					if s != "text" && s != "json" {
						return fmt.Errorf("log format must be text or json")
					}
					return nil
				},
			},
			&cli.StringFlag{
				Name:     flag_provider_name,
				Aliases:  []string{"n"},
				Usage:    "The name of the libdns provider",
				Sources:  cli.EnvVars(flagEnv(flag_provider_name)),
				Required: true,
				Validator: func(s string) error {
					if _, ok := libdns.Registry[s]; !ok {
						return fmt.Errorf("libdns provider %s not found", s)
					}
					return nil
				},
			},
			&cli.StringFlag{
				Name:    flag_provider_config,
				Aliases: []string{"c"},
				Usage:   "The json config for the libdns provider as a string",
				Sources: cli.EnvVars(flagEnv(flag_provider_config)),
			},
			&cli.StringSliceFlag{
				Name:      flag_provider_config_file,
				Aliases:   []string{"f"},
				TakesFile: true,
				Usage:     "The json config for the libdns provider read from `FILE`. Provide more than one file to merge configs.",
				Sources:   cli.EnvVars(flagEnv(flag_provider_config_file)),
			},
			&cli.StringFlag{
				Name:    flag_webhook_listen,
				Usage:   "The webhook server address to listen on (host:port)",
				Value:   "localhost:8888",
				Sources: cli.EnvVars(flagEnv(flag_webhook_listen)),
			},
			&cli.DurationFlag{
				Name:    flag_webhook_read_timeout,
				Usage:   "The webhook server read timeout",
				Value:   time.Minute,
				Sources: cli.EnvVars(flagEnv(flag_webhook_read_timeout)),
			},
			&cli.DurationFlag{
				Name:    flag_webhook_write_timeout,
				Usage:   "The webhook server write timeout",
				Value:   time.Minute,
				Sources: cli.EnvVars(flagEnv(flag_webhook_write_timeout)),
			},
			&cli.StringFlag{
				Name:    flag_metrics_listen,
				Usage:   "The metrics/health server address to listen on (host:port)",
				Value:   ":8080",
				Sources: cli.EnvVars(flagEnv(flag_metrics_listen)),
			},
			&cli.StringSliceFlag{
				Name:     flag_zones,
				Usage:    "The name of the dns zones that the webhook will manage",
				Required: true,
				Sources:  cli.EnvVars(flagEnv(flag_zones)),
			},
		},
		Action: func(ctx context.Context, cmd *cli.Command) error {
			if cmd.String(flag_log_format) == "json" {
				log.SetFormatter(&log.JSONFormatter{})
			} else {
				log.SetFormatter(&log.TextFormatter{})
			}

			level, err := log.ParseLevel(cmd.String(flag_log_level))
			if err != nil {
				return fmt.Errorf("failed to parse log level: %w", err)
			}
			log.Info("Setting log level to ", level)
			log.SetLevel(level)

			confs := [][]byte{}
			if conf := cmd.String(flag_provider_config); conf != "" {
				confs = append(confs, []byte(conf))
			}

			for _, file := range cmd.StringSlice(flag_provider_config_file) {
				log.Info("Reading provider config from file ", file)
				confFileBytes, err := os.ReadFile(file)
				if err != nil {
					return fmt.Errorf("failed to read file %s: %w", file, err)
				}
				confs = append(confs, confFileBytes)
			}
			if len(confs) == 0 {
				return fmt.Errorf("No provider config is provided")
			}

			providerName := cmd.String(flag_provider_name)

			provider, err := libdns.Registry[providerName](confs)
			if err != nil {
				return fmt.Errorf("failed to create provider %s: %w", providerName, err)
			}

			externaldnsProvider := webhook.NewWebhookProvider(cmd.StringSlice(flag_zones), provider)
			webhook.Run(
				externaldnsProvider,
				cmd.String(flag_webhook_listen),
				cmd.String(flag_metrics_listen),
				cmd.Duration(flag_webhook_read_timeout),
				cmd.Duration(flag_webhook_write_timeout),
			)
			return nil
		},
	}

	if err := cmd.Run(context.Background(), os.Args); err != nil {
		log.Fatal(err)
	}
}
