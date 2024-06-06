// Package loader implements methods for config load.
package loader

import (
	"fmt"

	"test_task/internal/config"

	"github.com/spf13/viper"
)

// LoadConfig loads and validates configuration to config.Config.
func LoadConfig(dir, filename string, conf *config.Config) error {
	viper.AddConfigPath(dir)
	viper.SetConfigName(filename)

	if err := viper.ReadInConfig(); err != nil {
		return fmt.Errorf("failed to read config: %w", err)
	}

	err := viper.Unmarshal(&conf)
	if err != nil {
		return fmt.Errorf("failed to unmarshal config file: %w", err)
	}

	return nil
}
