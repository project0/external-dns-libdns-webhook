package main

//go:generate go run ./internal/generator/generate.go code providers.json internal/libdnsregistry/registry.go
import (
	"context"
	"errors"
	"fmt"
	"os"
	"slices"
	"strings"
	"time"

	"github.com/project0/external-dns-libdns-webhook/internal/libdnsregistry"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/urfave/cli/v3"
)

var (
	// set by goreleaser.
	version = "unknown"
	commit  = "unknown"
	date    = "unknown"
)

const (
	envPrefix               = "LIBDNS_"
	flagLogLevel            = "log.level"
	flagLogFormat           = "log.format"
	flagProviderName        = "provider.name"
	flagProviderConfig      = "provider.config"
	flagProviderZones       = "provider.zones"
	flagWebhookListen       = "webhook.listen"
	flagWebhookReadTimeout  = "webhook.read-timeout"
	flagWebhookWriteTimeout = "webhook.write-timeout"
)

func flagEnv(name string) cli.ValueSourceChain {
	return cli.EnvVars(
		strings.ToUpper(envPrefix +
			strings.ReplaceAll(
				strings.ReplaceAll(name, "-", "_"),
				".", "_",
			),
		))
}

func flagAlternatives(name string) cli.ValueSourceChain {
	// IDEA: Add more sources here
	return cli.NewValueSourceChain(
		flagEnv(name).Chain...,
	)
}

const description = `A generic external-dns webhook provider supporting several libdns based DNS providers.

Supported Providers:
%s

Build: version=%s,commit=%s,date=%s
`

func main() {
	cmd := &cli.Command{
		Name:        "external-dns-webhook-libdns",
		Version:     version,
		Usage:       "Webhook for external-dns using libdns providers.",
		Description: fmt.Sprintf(description, strings.Join(libdnsregistry.List(), ", "), version, commit, date),

		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    flagLogLevel,
				Usage:   "The log level (trace, debug, info, warn, error, fatal, panic)",
				Value:   "info",
				Sources: flagAlternatives(flagLogLevel),
				Validator: func(s string) error {
					if _, err := zerolog.ParseLevel(s); err != nil {
						return fmt.Errorf("cannot set log level: %w", err)
					}

					return nil
				},
			},
			&cli.StringFlag{
				Name:    flagLogFormat,
				Usage:   "The log format (text, json)",
				Value:   "text",
				Sources: flagAlternatives(flagLogFormat),
				Validator: func(s string) error {
					if s != "text" && s != "json" {
						return errors.New("log format must be text or json")
					}

					return nil
				},
			},
			&cli.StringFlag{
				Name:     flagProviderName,
				Usage:    "The name of the libdns provider",
				Sources:  flagAlternatives(flagProviderName),
				Required: true,
				Validator: func(s string) error {
					if !slices.Contains(libdnsregistry.List(), s) {
						return fmt.Errorf("libdns provider %s not found", s)
					}

					return nil
				},
			},
			&cli.StringFlag{
				Name:    flagProviderConfig,
				Usage:   "The json config for the libdns provider as a string",
				Sources: flagAlternatives(flagProviderConfig),
				Value:   "{}",
			},
			&cli.StringSliceFlag{
				Name:     flagProviderZones,
				Usage:    "The name of the dns zones that the provider will manage",
				Required: true,
				Sources:  flagAlternatives(flagProviderZones),
			},
			&cli.StringFlag{
				Name:    flagWebhookListen,
				Usage:   "The webhook server address to listen on (host:port)",
				Value:   "localhost:8888",
				Sources: flagAlternatives(flagWebhookListen),
			},
			&cli.DurationFlag{
				Name:    flagWebhookReadTimeout,
				Usage:   "The webhook server read timeout",
				Value:   time.Minute,
				Sources: flagAlternatives(flagWebhookReadTimeout),
			},
			&cli.DurationFlag{
				Name:    flagWebhookWriteTimeout,
				Usage:   "The webhook server write timeout",
				Value:   time.Minute,
				Sources: flagAlternatives(flagWebhookWriteTimeout),
			},
		},

		Before: func(ctx context.Context, cmd *cli.Command) (context.Context, error) {
			// setup logging
			if cmd.String(flagLogFormat) == "text" {
				log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: time.RFC3339})
			}
			level, err := zerolog.ParseLevel(cmd.String(flagLogLevel))
			if err != nil {
				return ctx, fmt.Errorf("failed to parse log level: %w", err)
			}
			log.Info().Msgf("Setting log level to %s", level.String())
			zerolog.SetGlobalLevel(level)

			return ctx, nil
		},

		Action: func(_ context.Context, cmd *cli.Command) error {
			confs := [][]byte{}
			if conf := cmd.String(flagProviderConfig); conf != "" {
				confs = append(confs, []byte(conf))
			}

			if len(confs) == 0 {
				return errors.New("no provider config is provided")
			}

			providerName := cmd.String(flagProviderName)

			_, err := libdnsregistry.New(providerName, confs)
			if err != nil {
				return fmt.Errorf("failed to create provider %s: %w", providerName, err)
			}

			// externaldnsProvider := webhook.NewWebhookProvider(cmd.StringSlice(flag_zones), provider)
			// webhook.Run(
			// 	externaldnsProvider,
			// 	cmd.String(flag_webhook_listen),
			// 	cmd.String(flag_metrics_listen),
			// 	cmd.Duration(flag_webhook_read_timeout),
			// 	cmd.Duration(flag_webhook_write_timeout),
			// )
			return nil
		},
	}

	if err := cmd.Run(context.Background(), os.Args); err != nil {
		log.Fatal().Err(err).Msg("Failed to run")
	}
}
