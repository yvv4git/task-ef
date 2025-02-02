package repository_test

import (
	"database/sql"
	"fmt"
	"path"

	_ "github.com/lib/pq"

	"github.com/go-testfixtures/testfixtures/v3"
	"github.com/jackc/pgx/v5/pgxpool"
)

const (
	dbDialect   = "postgresql"
	dbDriver    = "postgres"
	fixturesDir = "internal/infrastructure/repository/fixtures"
)

func loadFixtures(dbPool *pgxpool.Pool, subDir string) error {
	connConfig := dbPool.Config().ConnConfig
	dsn := fmt.Sprintf(
		"postgres://%s:%s@%s:%d/%s?sslmode=disable",
		connConfig.User,
		connConfig.Password,
		connConfig.Host,
		connConfig.Port,
		connConfig.Database,
	)

	db, err := sql.Open(dbDriver, dsn)
	if err != nil {
		return err
	}

	fixtures, err := testfixtures.New(
		testfixtures.Database(db),
		testfixtures.Dialect(dbDialect),
		testfixtures.Directory(path.Join(fixturesDir, subDir)),
		testfixtures.DangerousSkipTestDatabaseCheck(),
	)
	if err != nil {
		return err
	}

	if err = fixtures.Load(); err != nil {
		return err
	}

	return nil
}
