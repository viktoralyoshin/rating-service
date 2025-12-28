package config

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func setEnv(t *testing.T, key, value string) {
	if err := os.Setenv(key, value); err != nil {
		t.Fatalf("error while setting env: %v", err)
	}
}

func TestLoad(t *testing.T) {
	setEnv(t, "KAFKA_ADDR", "localhost:9092")
	setEnv(t, "PORT", "8080")
	setEnv(t, "GAME_SERVICE_ADDR", "games:50051")

	cfg := Load()

	assert.Equal(t, "localhost:9092", cfg.KafkaAddr)
	assert.Equal(t, "8080", cfg.Port)
	assert.Equal(t, "games:50051", cfg.GameServiceAddr)
}
