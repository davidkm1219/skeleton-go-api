package config_test

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/twk/skeleton-go-api/internal/config"
)

func TestViper_BuildConfig(t *testing.T) {
	t.Parallel()

	type args struct {
		path string
	}

	type want struct {
		config *config.Config
		err    error
	}

	tests := map[string]struct {
		args args
		want want
	}{
		"valid config": {
			args: args{
				path: "test/config.yaml",
			},
			want: want{
				config: &config.Config{
					ConfigPath: "test/config.yaml",
					LogLevel:   "info",
					Stacktrace: true,
				},
			},
		},
		"invalid path": {
			args: args{
				path: "test/not-existing.yaml",
			},
			want: want{
				config: &config.Config{
					ConfigPath: "test/not-existing.yaml",
				},
			},
		},
		"invalid yaml": {
			args: args{
				path: "test/notyaml.yaml",
			},
			want: want{
				err: errors.New("error reading local config file"),
			},
		},
	}
	for name, tt := range tests {
		tt := tt

		t.Run(name, func(t *testing.T) {
			t.Parallel()

			v := config.NewViper()
			v.Viper.Set("config_path", tt.args.path)

			cfg, err := v.BuildConfig()
			if tt.want.err != nil {
				assert.ErrorContains(t, err, tt.want.err.Error())
				return
			}

			assert.Equal(t, tt.want.config, cfg)
		})
	}
}
