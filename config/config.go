package config

type Config struct {
	ConfigName string        `json:"config_name"`
	Gin        *GinEntity    `json:"gin"`
	Gorm       []GormEntity  `json:"gorm"`
	Redis      []RedisEntity `json:"redis"`
}

// GinEntity Gin Framework config
type GinEntity struct {
	Port       string `json:"port" binding:"required"`
	Enable     bool   `json:"enable" binding:"required"`
	Mode       string `json:"mode"`
	Middleware *struct {
		Cors *struct {
			Enabled          bool     `json:"enabled" binding:"required"`
			AllowOrigins     []string `json:"allow_origins"`
			AllowMethods     []string `json:"allow_methods"`
			AllowHeaders     []string `json:"allow_headers"`
			ExposeHeaders    []string `json:"expose_headers"`
			AllowCredentials bool     `json:"allow_credentials"`
			MaxAge           int      `json:"max_age"`
		} `json:"cors"`
		StaticFile *struct {
			Enable             bool   `json:"enable"`
			MaxMultipartMemory int64  `json:"max_multipart_memory"`
			UrlOne             string `json:"url_one"`
			UrlTwo             string `json:"url_two"`
		} `json:"static_file"`
	} `json:"middleware"`
}

// GormEntity Gorm ORM config
type GormEntity struct {
	Enable   bool `json:"enable"`
	Database []struct {
		Host     string   `json:"host"`
		Port     string   `json:"port"`
		Username string   `json:"username"`
		Password string   `json:"password"`
		Database string   `json:"database"`
		Param    []string `json:"param"`
	}
}

// RedisEntity Redis config
type RedisEntity struct {
	Enable   bool   `json:"enable"`
	Addr     string `json:"addr"`
	Password string `json:"password"`
	DB       int    `json:"db"`
	Protocol int    `json:"protocol"`
}
