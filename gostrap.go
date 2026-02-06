package gostrap

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"gostrap/config"
	"gostrap/provider"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

var (
	configName = "gostrap"
)

type Bootstrap struct {
	Config     *config.Config
	Gin        *gin.Engine
	Database   map[string]*gorm.DB
	Redis      map[int]*redis.Client
	HttpServer *http.Server
}

func LetsGo(filePath string) (*Bootstrap, error) {
	var (
		bootstrap = new(Bootstrap)
		cfg       *config.Config
		err       error
	)
	if cfg, err = openConfigFile(filePath); err != nil {
		return bootstrap, err
	}
	return startApplicationProcess(cfg)
}

func (boot *Bootstrap) Pulse() error {
	if boot.Gin != nil {
		srv := &http.Server{
			Addr:    fmt.Sprintf(":%d", boot.Config.Gin.Port),
			Handler: boot.Gin,
		}
		go func() {
			if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
				log.Fatalf("[Gostrap] listen: %s\n", err)
			}
		}()
		boot.HttpServer = srv
	}
	return nil
}
func (boot *Bootstrap) Wait() {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)
	<-quit
	log.Println("[Gostrap] Shutting down server...")
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	boot.ShutDown(ctx)
}
func (boot *Bootstrap) ShutDown(ctx context.Context) {
	if boot.Gin != nil && boot.HttpServer != nil {
		if err := boot.HttpServer.Shutdown(ctx); err != nil {
			log.Printf("[Gostrap] Server forced to shutdown: %v", err)
		}
	}
	if boot.Database != nil {
		for name, db := range boot.Database {
			sqlDB, err := db.DB()
			if err == nil {
				log.Printf("[Gostrap] Closing database: %s", name)
				if err := sqlDB.Close(); err != nil {
					log.Printf("[Gostrap] Error closing database: %v", err)
				}
			}
		}
	}
	if boot.Redis != nil {
		for name, redisClient := range boot.Redis {
			log.Printf("[Gostrap] Closing redis: %s", name)
			if err := redisClient.Close(); err != nil {
				log.Printf("[Gostrap] Error closing redis: %v", err)
			}
		}
	}
	log.Println("[Gostrap] Goodbye!")
}

func startApplicationProcess(cfg *config.Config) (*Bootstrap, error) {
	if cfg == nil {
		return nil, errors.New("[Gostrap] config is nil")
	}
	var bootstrap = &Bootstrap{Config: cfg}
	for _, pr := range provider.ProvidersPipeline {
		result, err := pr.Init(cfg)
		if err != nil {
			return nil, err
		}
		if result != nil {
			mapResultProviderToBootstrap(bootstrap, result)
		}
	}
	return bootstrap, nil
}
func openConfigFile(path string) (*config.Config, error) {
	var (
		file []byte
		err  error
		cfg  *config.Config
	)
	if file, err = os.ReadFile(path); err != nil {
		return cfg, err
	}
	if err = json.Unmarshal(file, &cfg); err != nil {
		return cfg, err
	}
	if cfg.ConfigName == configName {
		return cfg, nil
	}
	return nil, errors.New("[Gostrap] config name not match")
}
func mapResultProviderToBootstrap(b *Bootstrap, res any) {
	switch v := res.(type) {
	case *gin.Engine:
		b.Gin = v
	case map[string]*gorm.DB:
		b.Database = v
	case map[int]*redis.Client:
		b.Redis = v
	}
}
