package provider

import (
	"github.com/SexyBobRiK/gostrap/config"
	"github.com/gin-contrib/cors"

	"github.com/gin-gonic/gin"
)

type GinProvider struct{}

func (GinProvider) ProviderInit(cfgGin config.GinEntity) (*gin.Engine, error) {
	if cfgGin.Enable {
		if cfgGin.Mode != "" {
			gin.SetMode(cfgGin.Mode)
		} else {
			gin.SetMode(gin.DebugMode)
		}
		ginEngine := gin.Default()
		if cfgGin.Middleware != nil {
			if cfgGin.Middleware.Cors.Enabled && cfgGin.Middleware.Cors != nil {
				ginEngine.Use(cors.New(cors.Config{
					AllowCredentials: cfgGin.Middleware.Cors.AllowCredentials,
					AllowOrigins:     cfgGin.Middleware.Cors.AllowOrigins,
					AllowMethods:     cfgGin.Middleware.Cors.AllowMethods,
					AllowHeaders:     cfgGin.Middleware.Cors.AllowHeaders,
					ExposeHeaders:    cfgGin.Middleware.Cors.ExposeHeaders,
				}))
				if cfgGin.Middleware.StaticFile.Enable && cfgGin.Middleware.StaticFile != nil {
					ginEngine.Static(cfgGin.Middleware.StaticFile.UrlOne, cfgGin.Middleware.StaticFile.UrlTwo)
					ginEngine.MaxMultipartMemory = cfgGin.Middleware.StaticFile.MaxMultipartMemory
				}
			} else {
				ginEngine.Use(cors.Default())
			}
		}
		return ginEngine, nil
	}
	return nil, nil
}
