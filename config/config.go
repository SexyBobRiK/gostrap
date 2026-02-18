package config

type Config struct {
	ConfigName string           `json:"config_name" yaml:"configName" yml:"configName"`
	Gin        *GinEntity       `json:"gin" yaml:"gin" yml:"gin"`
	Database   []DatabaseEntity `json:"database" yaml:"database" yml:"database"`
	Redis      []RedisEntity    `json:"redis" yaml:"redis" yml:"redis"`
}

// GinEntity Gin Framework config
type GinEntity struct {
	Port       string `json:"port" binding:"required" yaml:"port" yml:"port"`
	Enable     bool   `json:"enable" binding:"required" yaml:"enable" yml:"enable"`
	Mode       string `json:"mode" yaml:"mode" yml:"mode"`
	Middleware *struct {
		Cors *struct {
			Enabled          bool     `json:"enabled" binding:"required" yaml:"enabled" yml:"enabled"`
			AllowOrigins     []string `json:"allow_origins" yaml:"allowOrigins" yml:"allowOrigins"`
			AllowMethods     []string `json:"allow_methods" yaml:"allowMethods" yml:"allowMethods"`
			AllowHeaders     []string `json:"allow_headers" yaml:"allowHeaders" yml:"allowHeaders"`
			ExposeHeaders    []string `json:"expose_headers" yaml:"exposeHeaders" yml:"exposeHeaders"`
			AllowCredentials bool     `json:"allow_credentials" yaml:"allowCredentials" yml:"allowCredentials"`
			MaxAge           int      `json:"max_age" yaml:"maxAge" yml:"maxAge"`
		} `json:"cors" yaml:"cors" yml:"cors"`
		StaticFile *struct {
			Enable             bool   `json:"enable" binding:"required" yaml:"enable" yml:"enable"`
			MaxMultipartMemory int64  `json:"max_multipart_memory" binding:"required" yaml:"maxMultipartMemory" yml:"maxMultipartMemory"`
			UrlOne             string `json:"url_one" binding:"required" yaml:"urlOne" yml:"urlOne"`
			UrlTwo             string `json:"url_two" binding:"required" yaml:"urlTwo" yml:"urlTwo"`
		} `json:"static_file" yaml:"staticFile" yml:"staticFile"`
	} `json:"middleware" yaml:"middleware" yml:"middleware"`
}

// DatabaseEntity Gorm ORM config
type DatabaseEntity struct {
	Enable   bool `json:"enable" binding:"required" yaml:"enable" yml:"enable"`
	Database []struct {
		Enable   bool     `json:"enable" binding:"required" yaml:"enable" yml:"enable"`
		Host     string   `json:"host" binding:"required" yaml:"host" yml:"host"`
		Port     string   `json:"port" binding:"required" yaml:"port" yml:"port"`
		Username string   `json:"username" binding:"required" yaml:"username" yml:"username"`
		Password string   `json:"password" binding:"required" yaml:"password" yml:"password"`
		Database *string  `json:"database" yaml:"database" yml:"database"`
		Param    []string `json:"param" yaml:"param" yml:"param"`
	} `json:"database" yaml:"database" yml:"database"`
}

// RedisEntity Redis config
type RedisEntity struct {
	Enable   bool   `json:"enable" binding:"required" yaml:"enable" yml:"enable"`
	Addr     string `json:"addr" binding:"required" yaml:"addr" yml:"addr"`
	Password string `json:"password" yaml:"password" yml:"password"`
	DB       int    `json:"db" yaml:"db" yml:"db"`
	Protocol int    `json:"protocol" yaml:"protocol" yml:"protocol"`
}
