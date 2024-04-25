package commands

import (
	"fmt"

	"github.com/spf13/cobra"
	"go.uber.org/zap"

	"github.com/twk/skeleton-go-api/internal/config"
	"github.com/twk/skeleton-go-api/internal/logger"
)

// NewPlaceholderCmd creates a new cobra command for the get command
func NewPlaceholderCmd(v *config.Viper, l *logger.Logger) *cobra.Command {
	b := []config.BindDetail{
		{Flag: config.FlagDetail{Name: "id", Shorthand: "i", Description: "placeholder flag option", DefaultValue: 1}, MapKey: "placeholder.id"},
	}

	cmd := &cobra.Command{
		Use:   "placeholder",
		Short: "placeholder for a command that does nothing",
		Long:  `This command does nothing.`,
		RunE: func(_ *cobra.Command, _ []string) error {
			return startPlaceholder(v, l)
		},
	}

	if err := v.SetFlagAndBind(cmd, b); err != nil {
		return nil
	}

	return cmd
}

func startPlaceholder(v *config.Viper, log *logger.Logger) error {
	cfg, err := v.BuildConfig()
	if err != nil {
		return fmt.Errorf("error building config: %w", err)
	}

	log.Info("starting", zap.Any("config", cfg))

	return nil
}
