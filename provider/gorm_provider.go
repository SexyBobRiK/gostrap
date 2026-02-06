package provider

import (
	"fmt"
	"gostrap/config"
	"log"
	"strings"
	"sync"

	"golang.org/x/sync/errgroup"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type GormProvider struct{}

func (GormProvider) ProviderInit(entities []config.GormEntity) (map[string]*gorm.DB, error) {
	var (
		dbMap = make(map[string]*gorm.DB)
		mu    sync.Mutex
		eg    errgroup.Group
	)

	for _, entity := range entities {
		if !entity.Enable {
			continue
		}
		for _, dbCfg := range entity.Database {
			eg.Go(func() error {
				dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s %s",
					dbCfg.Host,
					dbCfg.Username,
					dbCfg.Password,
					dbCfg.Database,
					dbCfg.Port,
					strings.Join(dbCfg.Param, " "),
				)
				db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
				if err != nil {
					return fmt.Errorf("[Gostrap] Database %s connect error: %w", dbCfg.Database, err)
				}
				mu.Lock()
				dbMap[dbCfg.Database] = db
				mu.Unlock()
				log.Printf("[Gostrap] Database %s connected", dbCfg.Database)
				return nil
			})
		}
	}
	if err := eg.Wait(); err != nil {
		return nil, err
	}
	return dbMap, nil
}
