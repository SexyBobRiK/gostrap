package provider

import (
	"strings"

	"github.com/SexyBobRiK/gostrap/config"
	"github.com/gin-contrib/cors"

	"github.com/gin-gonic/gin"
)

type GinProvider struct{}

func (GinProvider) ProviderInit(cfgGin config.GinEntity) (*gin.Engine, error) {
	if cfgGin.Enable {
		if cfgGin.Mode != "" && (strings.ToLower(cfgGin.Mode) == gin.DebugMode ||
			strings.ToLower(cfgGin.Mode) == gin.ReleaseMode ||
			strings.ToLower(cfgGin.Mode) == gin.TestMode) {
			gin.SetMode(strings.ToLower(cfgGin.Mode))
		} else {
			gin.SetMode(gin.DebugMode)
		}
		ginEngine := gin.Default()
		if cfgGin.Middleware != nil {
			if cfgGin.Middleware.Cors != nil && cfgGin.Middleware.Cors.Enabled {
				if cfgGin.Middleware.Cors.AllowOrigins == nil || len(cfgGin.Middleware.Cors.AllowOrigins) == 0 {
					cfgGin.Middleware.Cors.AllowOrigins = []string{"*"}
				}
				ginEngine.Use(cors.New(cors.Config{
					AllowCredentials: cfgGin.Middleware.Cors.AllowCredentials,
					AllowOrigins:     cfgGin.Middleware.Cors.AllowOrigins,
					AllowMethods:     cfgGin.Middleware.Cors.AllowMethods,
					AllowHeaders:     cfgGin.Middleware.Cors.AllowHeaders,
					ExposeHeaders:    cfgGin.Middleware.Cors.ExposeHeaders,
				}))
			} else {
				ginEngine.Use(cors.Default())
			}
			if cfgGin.Middleware.StaticFile != nil && cfgGin.Middleware.StaticFile.Enable {
				ginEngine.Static(cfgGin.Middleware.StaticFile.UrlOne, cfgGin.Middleware.StaticFile.UrlTwo)
				ginEngine.MaxMultipartMemory = cfgGin.Middleware.StaticFile.MaxMultipartMemory
			}
		}
		return ginEngine, nil
	}
	return nil, nil
}
