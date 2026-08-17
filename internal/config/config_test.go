package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLoadReadsSchedulerAndRetrySettings(t *testing.T) {
	path := filepath.Join(t.TempDir(), "config.yaml")
	data := []byte("server:\n  scheduler:\n    host: 127.0.0.1\n    port: 8080\n    grpc_port: 9090\nretry:\n  max_attempts: 3\n  initial_delay: 1s\n  max_delay: 2s\n")
	if err := os.WriteFile(path, data, 0600); err != nil {
		t.Fatal(err)
	}
	cfg, err := Load(path)
	if err != nil {
		t.Fatal(err)
	}
	if cfg.Server.Scheduler.GRPCPort != 9090 || cfg.Retry.MaxAttempts != 3 {
		t.Fatalf("unexpected config: %#v", cfg)
	}
}
