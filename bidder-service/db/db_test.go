package db

import (
	"os"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

func TestMigrate(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	// Use temporary migration folder
	//os.Setenv("MIGRATIONS_DIR", "./test_migrations")

	// goose expects a real *sql.DB with a registered driver name
	//	sqlDB := goose.WrapDB(db)

	// --- Expected Goose behavior ---
	mock.ExpectBegin()

	// goose checks current version
	mock.ExpectQuery("(?i)^SELECT version_id, is_applied FROM goose_db_version ORDER BY id DESC$").
		WillReturnRows(sqlmock.NewRows([]string{"version_id", "is_applied"}))

	// goose tries to create goose table if not exists
	mock.ExpectExec("(?i)^CREATE TABLE IF NOT EXISTS goose_db_version").
		WillReturnResult(sqlmock.NewResult(0, 0))

	mock.ExpectCommit()

	// --- Run migration ---
	err = migrate(db)
	assert.NoError(t, err)

	// Make sure all expectations are met
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestMain(m *testing.M) {
	// --- Pre-testdata setup ---
	abs, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	err = os.Setenv("MIGRATIONS_DIR", abs+"/migrations")

	if err != nil {
		panic(err)
	}

	// --- Run tests ---
	code := m.Run()

	os.Exit(code)
}
