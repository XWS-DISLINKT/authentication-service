package config

type Config struct {
	Port        string
	ProfilePort string
	ProfileHost string
}

func NewConfig() *Config {
	return &Config{
		Port:        "8003",
		ProfilePort: "8001",
		ProfileHost: "localhost",
	}
}
