package config

import (
	"encoding/json"
	"errors"
	"path/filepath"

	"github.com/goccy/go-yaml"
)

type LoaderConfig interface {
	LoadConfig(file []byte, config *Config) error
}
type JSONLoaderConfig struct{}

func (JSONLoaderConfig) LoadConfig(file []byte, config *Config) error {
	return json.Unmarshal(file, config)
}

type YAMLLoaderConfig struct{}

func (YAMLLoaderConfig) LoadConfig(file []byte, config *Config) error {
	return yaml.Unmarshal(file, config)
}

var loaders = map[string]LoaderConfig{
	".json": JSONLoaderConfig{},
	".yaml": YAMLLoaderConfig{},
	".yml":  YAMLLoaderConfig{},
}

func Decoder(path string) (LoaderConfig, error) {
	ext := filepath.Ext(path)
	loader, ok := loaders[ext]
	if !ok {
		return nil, errors.New("[Gostrap] config file must be .json or .yaml or .yml")
	}
	return loader, nil
}
