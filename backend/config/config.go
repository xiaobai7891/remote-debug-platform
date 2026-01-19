package config

import (
	"os"
	"sync"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Server   ServerConfig   `yaml:"server"`
	Database DatabaseConfig `yaml:"database"`
	PortPool PortPoolConfig `yaml:"port_pool"`
	JWT      JWTConfig      `yaml:"jwt"`
	Pricing  PricingConfig  `yaml:"pricing"`
	Admin    AdminConfig    `yaml:"admin"`
}

type ServerConfig struct {
	HTTPPort int `yaml:"http_port"`
	WSPort   int `yaml:"ws_port"`
}

type DatabaseConfig struct {
	Path string `yaml:"path"`
}

type PortPoolConfig struct {
	Min int `yaml:"min"`
	Max int `yaml:"max"`
}

type JWTConfig struct {
	Secret      string `yaml:"secret"`
	ExpireHours int    `yaml:"expire_hours"`
}

type PricingConfig struct {
	PortPerDay float64 `yaml:"port_per_day"`
}

type AdminConfig struct {
	Username string `yaml:"username"`
	Password string `yaml:"password"`
}

var (
	cfg  *Config
	once sync.Once
)

func Load(path string) (*Config, error) {
	var err error
	once.Do(func() {
		cfg = &Config{}
		data, readErr := os.ReadFile(path)
		if readErr != nil {
			err = readErr
			return
		}
		err = yaml.Unmarshal(data, cfg)
	})
	return cfg, err
}

func Get() *Config {
	if cfg == nil {
		// 返回默认配置
		return &Config{
			Server: ServerConfig{
				HTTPPort: 3000,
				WSPort:   3001,
			},
			Database: DatabaseConfig{
				Path: "./data/app.db",
			},
			PortPool: PortPoolConfig{
				Min: 10000,
				Max: 60000,
			},
			JWT: JWTConfig{
				Secret:      "remote-debug-secret-key-2024",
				ExpireHours: 168,
			},
			Pricing: PricingConfig{
				PortPerDay: 1.0,
			},
			Admin: AdminConfig{
				Username: "admin",
				Password: "admin123",
			},
		}
	}
	return cfg
}

func SetDefault() {
	cfg = &Config{
		Server: ServerConfig{
			HTTPPort: 3000,
			WSPort:   3001,
		},
		Database: DatabaseConfig{
			Path: "./data/app.db",
		},
		PortPool: PortPoolConfig{
			Min: 10000,
			Max: 60000,
		},
		JWT: JWTConfig{
			Secret:      "remote-debug-secret-key-2024",
			ExpireHours: 168,
		},
		Pricing: PricingConfig{
			PortPerDay: 1.0,
		},
		Admin: AdminConfig{
			Username: "admin",
			Password: "admin123",
		},
	}
}
