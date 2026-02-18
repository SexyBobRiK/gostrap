package provider

import (
	"fmt"
	"log"
	"strings"
	"sync"

	"github.com/SexyBobRiK/gostrap/config"

	"golang.org/x/sync/errgroup"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type GormProvider struct{}

func (GormProvider) ProviderInit(entities []config.DatabaseEntity) (map[string]gorm.DB, error) {
	var (
		dbMap               = make(map[string]gorm.DB)
		mu                  sync.Mutex
		eg                  errgroup.Group
		defaultDatabaseName = "postgres"
	)

	for _, entity := range entities {
		if !entity.Enable {
			continue
		}
		for i, _ := range entity.Database {
			if !entity.Database[i].Enable {
				continue
			}
			database := entity.Database[i]
			eg.Go(func() error {
				if database.Database == nil || *database.Database == "" {
					database.Database = &defaultDatabaseName
				}
				dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s %s",
					database.Host,
					database.Username,
					database.Password,
					database.Database,
					database.Port,
					strings.Join(database.Param, " "),
				)
				db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
				if err != nil {
					return fmt.Errorf("[Gostrap] Database %s connect error: %w", database.Database, err)
				}
				mu.Lock()
				dbMap[*database.Database] = *db
				mu.Unlock()
				log.Printf("[Gostrap] Database %s connected", database.Database)
				return nil
			})
		}
	}
	if err := eg.Wait(); err != nil {
		return nil, err
	}
	return dbMap, nil
}
