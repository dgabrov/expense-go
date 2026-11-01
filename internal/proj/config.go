package proj

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/dgabrov/expense-go/internal/data"
)

const configEnvVar = "CONFIG_FILE"

func loadConfig() (data.ConfigData, error) {
	var cfg data.ConfigData

	path := os.Getenv(configEnvVar)
	if path == "" {
		return cfg, fmt.Errorf("environment variable %s is not set", configEnvVar)
	}

	raw, err := os.ReadFile(path)
	if err != nil {
		return cfg, fmt.Errorf("reading config file %s: %w", path, err)
	}

	if err := json.Unmarshal(raw, &cfg); err != nil {
		return cfg, fmt.Errorf("parsing config file %s: %w", path, err)
	}

	return cfg, nil
}
