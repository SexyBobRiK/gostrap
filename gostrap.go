package gostrap

import (
	"context"
	"errors"
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/SexyBobRiK/gostrap/config"
	"github.com/SexyBobRiK/gostrap/provider"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

const (
	configName = "gostrap"
)

type Bootstrap struct {
	Logger     *slog.Logger
	Config     *config.Config
	Gin        *gin.Engine
	Database   map[string]gorm.DB
	Redis      map[int]redis.Client
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

func (b *Bootstrap) initLogger() {
	var handler slog.Handler
	if strings.ToLower(b.Config.Gin.Mode) == gin.DebugMode || strings.ToLower(b.Config.Gin.Mode) == gin.TestMode {
		handler = slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug})
	} else {
		handler = slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo})
	}
	b.Logger = slog.New(handler)
}
func (boot *Bootstrap) Pulse() error {
	if boot.Gin != nil {
		srv := &http.Server{
			Addr:    fmt.Sprintf(":%s", boot.Config.Gin.Port),
			Handler: boot.Gin,
		}
		go func() {
			if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
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
			log.Printf("[Gostrap] Closing redis: %d", name)
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
	bootstrap.initLogger()
	return bootstrap, nil
}
func openConfigFile(path string) (*config.Config, error) {
	var (
		file []byte
		err  error
		cfg  = new(config.Config)
	)
	if file, err = os.ReadFile(path); err != nil {
		return nil, err
	}
	loader, err := config.Decoder(path)
	if err != nil {
		return nil, err
	}
	err = loader.LoadConfig(file, cfg)
	if err != nil {
		return nil, err
	}
	if cfg.ConfigName != configName {
		return nil, errors.New("[Gostrap] config name not match")
	}
	return cfg, nil
}
func mapResultProviderToBootstrap(b *Bootstrap, res any) {
	switch v := res.(type) {
	case *gin.Engine:
		b.Gin = v
	case map[string]gorm.DB:
		b.Database = v
	case map[int]redis.Client:
		b.Redis = v
	}
}
