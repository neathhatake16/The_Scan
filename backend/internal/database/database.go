package database

import (
	"fmt"
	"log"
	"path/filepath"
	"time"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/neathhatake/the_Scan/internal/config"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func ConnectDatabase(cfg *config.Config) *gorm.DB {
	var db *gorm.DB
	var err error

	for i := 0; i < 10; i++ {
		db, err = gorm.Open(postgres.Open(cfg.GetDSN()), &gorm.Config{
			Logger: logger.Default.LogMode(logger.Info),
		})
		if err == nil {
			break
		}
		log.Printf("DB not ready yet (%d/10), retrying in 3s... ⏳", i+1)
		time.Sleep(3 * time.Second)

	}
	if err != nil {
		log.Fatalf("Failed to connect to database: %v  ❌ ", err)
	}

	sqlDB, _ := db.DB()
	sqlDB.SetMaxOpenConns(25)
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetConnMaxLifetime(5 * time.Minute)

	log.Println("Database connected .....✅")
	return db

}

func RunMigrations(db *gorm.DB, migrationsPath string, dsn string) {
	absPath, err := filepath.Abs(migrationsPath)
	if err != nil {
		log.Fatalf("Failed to resolve migrations path: %v", err)
	}

	m, err := migrate.New(
		fmt.Sprintf("file://%s", absPath),
		dsn,
	)
	if err != nil {
		log.Fatalf("Failed to create migration instance: %v", err)
	}

	err = m.Up()
	if err != nil && err != migrate.ErrNoChange {
		log.Fatalf("Failed to run migrations: %v", err)
	}

	if err == migrate.ErrNoChange {
		log.Println("No new migrations to apply")
	} else {
		log.Println("Migrations applied successfully")
	}
}

