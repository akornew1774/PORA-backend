// database - пакет для подключения БД к проекту
package database

import (
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	"gorm.io/gorm"
)

// RunMigrations применяет существующие миграции к БД
func RunMigrations(db *gorm.DB) error {
	DB, err := db.DB()
	if err != nil {
		return err
	}

	driver, err := postgres.WithInstance(DB, &postgres.Config{})
	if err != nil {
		return err
	}

	migrations, err := migrate.NewWithDatabaseInstance(
		"file://migrations",
		"postgres",
		driver,
	)

	if err != nil {
		return err
	}

	err = migrations.Up()
	if err != nil && err != migrate.ErrNoChange {
		return err
	}

	return nil
}
