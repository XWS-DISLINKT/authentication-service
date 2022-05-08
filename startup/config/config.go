package config

import (
	"fmt"
	"os"
)

type Config struct {
	Port        string
	ProfilePort string
	ProfileHost string
}

func NewConfig() *Config {
	if _, err := os.Stat("/.dockerenv"); err == nil {
		fmt.Println("docker")

		return &Config{
			Port:        os.Getenv("AUTHENTICATION_SERVICE_PORT"),
			ProfilePort: os.Getenv("PROFILE_SERVICE_PORT"),
			ProfileHost: os.Getenv("PROFILE_SERVICE_HOST"),
		}
	} else {
		fmt.Println("local")

		return &Config{
			Port:        "8003",
			ProfilePort: "8001",
			ProfileHost: "localhost",
		}
	}
}
