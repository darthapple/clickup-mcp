package config

import (
	"testing"
	"time"
)

func TestLoadRequiresToken(t *testing.T) {
	t.Setenv("CLICKUP_API_TOKEN", "")
	t.Setenv("CLICKUP_TEAM_ID", "123")

	if _, err := Load(); err == nil {
		t.Fatal("expected error when CLICKUP_API_TOKEN is unset")
	}
}

func TestLoadRequiresTeamID(t *testing.T) {
	t.Setenv("CLICKUP_API_TOKEN", "pk_abc")
	t.Setenv("CLICKUP_TEAM_ID", "")

	if _, err := Load(); err == nil {
		t.Fatal("expected error when CLICKUP_TEAM_ID is unset")
	}
}

func TestLoadDefaults(t *testing.T) {
	t.Setenv("CLICKUP_API_TOKEN", "pk_abc")
	t.Setenv("CLICKUP_TEAM_ID", "123")
	t.Setenv("CLICKUP_API_BASE_URL", "")
	t.Setenv("CLICKUP_API_BASE_URL_V3", "")
	t.Setenv("CLICKUP_HTTP_TIMEOUT", "")
	t.Setenv("CLICKUP_MAX_RETRIES", "")

	cfg, err := Load()
	if err != nil {
		t.Fatalf("Load: %v", err)
	}
	if cfg.BaseURLv2 != "https://api.clickup.com/api/v2" {
		t.Errorf("BaseURLv2 = %q", cfg.BaseURLv2)
	}
	if cfg.BaseURLv3 != "https://api.clickup.com/api/v3" {
		t.Errorf("BaseURLv3 = %q", cfg.BaseURLv3)
	}
	if cfg.HTTPTimeout != 30*time.Second {
		t.Errorf("HTTPTimeout = %v", cfg.HTTPTimeout)
	}
	if cfg.MaxRetries != 4 {
		t.Errorf("MaxRetries = %d", cfg.MaxRetries)
	}
}

func TestLoadOverrides(t *testing.T) {
	t.Setenv("CLICKUP_API_TOKEN", "pk_abc")
	t.Setenv("CLICKUP_TEAM_ID", "123")
	t.Setenv("CLICKUP_HTTP_TIMEOUT", "10s")
	t.Setenv("CLICKUP_MAX_RETRIES", "7")

	cfg, err := Load()
	if err != nil {
		t.Fatalf("Load: %v", err)
	}
	if cfg.HTTPTimeout != 10*time.Second {
		t.Errorf("HTTPTimeout = %v", cfg.HTTPTimeout)
	}
	if cfg.MaxRetries != 7 {
		t.Errorf("MaxRetries = %d", cfg.MaxRetries)
	}
}

func TestLoadRejectsInvalidOverrides(t *testing.T) {
	t.Setenv("CLICKUP_API_TOKEN", "pk_abc")
	t.Setenv("CLICKUP_TEAM_ID", "123")
	t.Setenv("CLICKUP_MAX_RETRIES", "not-a-number")

	if _, err := Load(); err == nil {
		t.Fatal("expected error for invalid CLICKUP_MAX_RETRIES")
	}
}
