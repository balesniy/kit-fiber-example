package config

import (
	"os"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Server struct {
		Port            string `yaml:"port"`
		ShutdownTimeout int    `yaml:"shutdownTimeout"`
	} `yaml:"server"`
	RateLimit struct {
		Requests int    `yaml:"requests"`
		Duration string `yaml:"duration"`
	} `yaml:"rateLimit"`
	CircuitBreaker struct {
		Threshold   int    `yaml:"threshold"`
		Timeout     string `yaml:"timeout"`
		MaxRequests int    `yaml:"maxRequests"`
	} `yaml:"circuitBreaker"`
	Telemetry struct {
		ServiceName   string  `yaml:"serviceName"`
		CollectorAddr string  `yaml:"collectorAddr"`
		SamplingRatio float64 `yaml:"samplingRatio"`
	} `yaml:"telemetry"`
	Claude struct {
		APIKey     string `yaml:"apiKey"`
		BaseURL    string `yaml:"baseURL"`
		Model      string `yaml:"model"`
		Timeout    int    `yaml:"timeout"`
		MaxRetries int    `yaml:"maxRetries"`
	} `yaml:"claude"`
}

func LoadConfig(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var config Config
	if err := yaml.Unmarshal(data, &config); err != nil {
		return nil, err
	}

	return &config, nil
}
