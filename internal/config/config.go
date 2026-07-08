// Package config loads the ClickUp MCP server's configuration from
// environment variables.
package config

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

type Config struct {
	APIToken    string
	TeamID      string
	BaseURLv2   string
	BaseURLv3   string
	HTTPTimeout time.Duration
	MaxRetries  int
}

// Load reads and validates configuration from the environment.
// CLICKUP_API_TOKEN and CLICKUP_TEAM_ID are required; everything else has a
// default.
func Load() (*Config, error) {
	token := strings.TrimSpace(os.Getenv("CLICKUP_API_TOKEN"))
	if token == "" {
		return nil, fmt.Errorf("CLICKUP_API_TOKEN is required but not set")
	}
	if !strings.HasPrefix(token, "pk_") {
		fmt.Fprintln(os.Stderr, "clickup-mcp: warning: CLICKUP_API_TOKEN does not look like a personal token (expected pk_ prefix)")
	}

	teamID := strings.TrimSpace(os.Getenv("CLICKUP_TEAM_ID"))
	if teamID == "" {
		return nil, fmt.Errorf("CLICKUP_TEAM_ID is required but not set")
	}

	httpTimeout, err := envDuration("CLICKUP_HTTP_TIMEOUT", 30*time.Second)
	if err != nil {
		return nil, err
	}
	maxRetries, err := envInt("CLICKUP_MAX_RETRIES", 4)
	if err != nil {
		return nil, err
	}

	return &Config{
		APIToken:    token,
		TeamID:      teamID,
		BaseURLv2:   envOrDefault("CLICKUP_API_BASE_URL", "https://api.clickup.com/api/v2"),
		BaseURLv3:   envOrDefault("CLICKUP_API_BASE_URL_V3", "https://api.clickup.com/api/v3"),
		HTTPTimeout: httpTimeout,
		MaxRetries:  maxRetries,
	}, nil
}

func envOrDefault(key, def string) string {
	if v := strings.TrimSpace(os.Getenv(key)); v != "" {
		return v
	}
	return def
}

func envDuration(key string, def time.Duration) (time.Duration, error) {
	v := strings.TrimSpace(os.Getenv(key))
	if v == "" {
		return def, nil
	}
	d, err := time.ParseDuration(v)
	if err != nil {
		return 0, fmt.Errorf("invalid %s: %w", key, err)
	}
	return d, nil
}

func envInt(key string, def int) (int, error) {
	v := strings.TrimSpace(os.Getenv(key))
	if v == "" {
		return def, nil
	}
	n, err := strconv.Atoi(v)
	if err != nil {
		return 0, fmt.Errorf("invalid %s: %w", key, err)
	}
	return n, nil
}
