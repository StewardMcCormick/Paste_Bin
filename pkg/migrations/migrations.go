package migrations

import (
	"fmt"
	"github.com/golang-migrate/migrate/v4"

	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func Exec(dbUrl string, migrationsPath string) error {
	m, err := migrate.New(
		fmt.Sprintf("file://%s", migrationsPath),
		dbUrl,
	)
	if err != nil {
		return err
	}
	defer m.Close()

	if err = m.Up(); err != nil {
		return err
	}

	return nil
}
