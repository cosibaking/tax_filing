package ocrspace

import (
	"fmt"
	"os"
	"strings"
	"time"

	"gopkg.in/yaml.v3"
)

const (
	envAPIKey      = "OCRSPACE_API_KEY"
	envConfigPath  = "OCRSPACE_CONFIG"
	defaultTimeout = 60 * time.Second
)

// Config OCR.space 客户端配置
type Config struct {
	APIKey          string        `yaml:"apiKey"`
	Endpoint        string        `yaml:"endpoint"`
	EndpointGet     string        `yaml:"endpointGet"`
	Timeout         time.Duration `yaml:"timeout"`
	DefaultLanguage Language      `yaml:"defaultLanguage"`
	DefaultEngine   OCREngine     `yaml:"defaultEngine"`
}

type configFile struct {
	OCRSpace Config `yaml:"ocrspace"`
}

// DefaultConfig 返回带默认值的配置（不含 API Key）
func DefaultConfig() Config {
	return Config{
		Endpoint:        defaultPostEndpoint,
		EndpointGet:     defaultGetEndpoint,
		Timeout:         defaultTimeout,
		DefaultLanguage: LanguageEnglish,
		DefaultEngine:   Engine2,
	}
}

// LoadConfig 从 YAML 配置文件加载，路径为空时使用 OCRSPACE_CONFIG 或 manifest/config/ocrspace.yaml
func LoadConfig(path string) (Config, error) {
	cfg := DefaultConfig()

	if path == "" {
		path = os.Getenv(envConfigPath)
	}
	if path == "" {
		path = "manifest/config/ocrspace.yaml"
	}

	data, err := os.ReadFile(path)
	if err != nil {
		return cfg, fmt.Errorf("ocrspace: read config %q: %w", path, err)
	}
	return LoadConfigFromBytes(data)
}

// LoadConfigFromBytes 从 YAML 字节加载配置
func LoadConfigFromBytes(data []byte) (Config, error) {
	cfg := DefaultConfig()

	var file configFile
	if err := yaml.Unmarshal(data, &file); err != nil {
		return cfg, fmt.Errorf("ocrspace: parse config: %w", err)
	}

	mergeConfig(&cfg, file.OCRSpace)

	if key := strings.TrimSpace(os.Getenv(envAPIKey)); key != "" {
		cfg.APIKey = key
	}

	if err := cfg.Validate(); err != nil {
		return cfg, err
	}
	return cfg, nil
}

func mergeConfig(dst *Config, src Config) {
	if src.APIKey != "" {
		dst.APIKey = src.APIKey
	}
	if src.Endpoint != "" {
		dst.Endpoint = src.Endpoint
	}
	if src.EndpointGet != "" {
		dst.EndpointGet = src.EndpointGet
	}
	if src.Timeout > 0 {
		dst.Timeout = src.Timeout
	}
	if src.DefaultLanguage != "" {
		dst.DefaultLanguage = src.DefaultLanguage
	}
	if src.DefaultEngine > 0 {
		dst.DefaultEngine = src.DefaultEngine
	}
}

// Validate 校验配置
func (c Config) Validate() error {
	if strings.TrimSpace(c.APIKey) == "" {
		return fmt.Errorf("ocrspace: apiKey is required (set in config or %s env)", envAPIKey)
	}
	if c.Endpoint == "" {
		return fmt.Errorf("ocrspace: endpoint is required")
	}
	return nil
}
