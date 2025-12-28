package config

import "os"

type Config struct {
	KafkaAddr         string
	GameServiceAddr   string
	SocialServiceAddr string
	Env               string
	Port              string
}

func Load() *Config {
	return &Config{
		KafkaAddr:         os.Getenv("KAFKA_ADDR"),
		Env:               os.Getenv("ENV"),
		Port:              os.Getenv("PORT"),
		GameServiceAddr:   os.Getenv("GAME_SERVICE_ADDR"),
		SocialServiceAddr: os.Getenv("SOCIAL_SERVICE_ADDR"),
	}
}
