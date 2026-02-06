package provider

import "gostrap/config"

type ProviderInitiator[T any, R any] interface {
	ProviderInit(cfg T) (R, error)
}
type ProviderS struct {
	Name string
	Init func(config *config.Config) (any, error)
}

var ProvidersPipeline = []ProviderS{
	{
		Name: "gorm",
		Init: func(cfg *config.Config) (any, error) {
			p := GormProvider{}
			return p.ProviderInit(cfg.Gorm)
		},
	},
	{
		Name: "redis",
		Init: func(cfg *config.Config) (any, error) {
			p := RedisProvider{}
			return p.ProviderInit(cfg.Redis)
		},
	},
	{
		Name: "gin",
		Init: func(cfg *config.Config) (any, error) {
			if cfg.Gin != nil {
				p := GinProvider{}
				return p.ProviderInit(*cfg.Gin)
			}
			return nil, nil
		},
	},
}
