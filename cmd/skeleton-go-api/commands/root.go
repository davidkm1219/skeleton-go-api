// Package commands provides the command line interface for the application. It contains the root command and all the subcommands.
package commands

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/spf13/cobra"
	"go.uber.org/zap"

	"github.com/twk/skeleton-go-api/internal/api"
	"github.com/twk/skeleton-go-api/internal/client"
	"github.com/twk/skeleton-go-api/internal/config"
	"github.com/twk/skeleton-go-api/internal/logger"
	"github.com/twk/skeleton-go-api/internal/photos"
	"github.com/twk/skeleton-go-api/internal/server"
)

const appName = "skeleton-go-api"

// NewRootCommand creates a new cobra command for the root command
func NewRootCommand(l *logger.Logger) (*cobra.Command, error) {
	v := config.NewViper()

	b := []config.BindDetail{
		{Flag: config.FlagDetail{Name: "config", Description: fmt.Sprintf("Specifies the path to the configuration file for %s.", appName), DefaultValue: "./config.yaml"}, MapKey: "config_path"},
		{Flag: config.FlagDetail{Name: "log-level", Description: "Determines the logging verbosity level for the application. Available options are 'debug', 'info', 'warn', and 'error'.", DefaultValue: ""}, EnvName: "LOG_LEVEL", MapKey: "log_level"},
		{Flag: config.FlagDetail{Name: "stacktrace", Description: "Enables or disables the inclusion of stack traces in the log output.", DefaultValue: false}, EnvName: "STACKTRACE", MapKey: "stacktrace"},
	}

	rootCmd := &cobra.Command{
		Use:   appName,
		Short: "CLI for the skeleton-go-api application",
		Long: `CLI for the skeleton-go-api application.
This CLI is used to interact with the skeleton-go-api application.`,
		RunE: func(_ *cobra.Command, _ []string) error {
			return startRoot(v, l)
		},
		SilenceUsage: true,
	}

	if err := v.SetFlagAndBind(rootCmd, b); err != nil {
		return nil, fmt.Errorf("error initializing flags: %w", err)
	}

	rootCmd.AddCommand(NewPlaceholderCmd(v, l))

	return rootCmd, nil
}

func startRoot(v *config.Viper, l *logger.Logger) error {
	cfg, err := v.BuildConfig()
	if err != nil {
		return fmt.Errorf("error building config: %w", err)
	}

	l.Info("starting", zap.Any("config", cfg))

	httpClient := &http.Client{}
	hc, err := client.NewClient(httpClient)
	if err != nil {
		return fmt.Errorf("error creating http client: %w", err)
	}
	ps := photos.NewService(hc, l)
	pr := api.Photos(&cfg.Server, ps, l)
	rp := []server.RouteParam{
		{Method: http.MethodGet, Path: "/photos/:id", Handler: pr},
	}
	s := server.NewServer(&cfg.Server, gin.Default(), rp, l)

	if err := s.Start(); err != nil {
		return fmt.Errorf("error starting server: %w", err)
	}

	return nil
}
