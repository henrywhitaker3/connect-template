package postgres

import (
	"database/sql"
	"fmt"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	"github.com/golang-migrate/migrate/v4/source/iofs"
	"github.com/henrywhitaker3/connect-template/database/migrations"
)

type Migrator struct {
	m *migrate.Migrate
}

func NewMigrator(db *sql.DB) (*Migrator, error) {
	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		return nil, err
	}
	fs, err := iofs.New(migrations.Migrations, "files")
	if err != nil {
		return nil, fmt.Errorf("invalid migration fs: %w", err)
	}
	m, err := migrate.NewWithInstance(
		"iofs",
		fs,
		"postgres",
		driver,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create migrator: %w", err)
	}
	return &Migrator{
		m: m,
	}, nil
}

func (m *Migrator) Up() error {
	err := m.m.Up()
	if err == nil || err.Error() == "no change" || err.Error() == "first files: file does not exist" {
		return nil
	}
	return err
}

func (m *Migrator) Down() error {
	err := m.m.Down()
	if err == nil || err.Error() == "no change" {
		return nil
	}
	return err
}
