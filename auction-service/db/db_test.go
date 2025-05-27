package db

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/pressly/goose/v3"
)

func Init(m *testing.M) {
	tempDir, err := os.MkdirTemp("", "migrations_test")
	if err != nil {
		panic("failed to create temp dir for migrations")
	}
	os.Setenv("MIGRATIONS_DIR", tempDir)
}

func TestMigrateSuccess(t *testing.T) {
	mockDB, _, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed initating test %v", err)
	}
	// Enable logging to track actual SQL queries
	goose.SetLogger(log.New(os.Stdout, "[GOOSE] ", log.LstdFlags))

	// Override gooseUp to simulate success
	originalGooseUp := gooseUp
	gooseUp = func(db *sql.DB, dir string, optionsFunc ...goose.OptionsFunc) error {
		return nil
	}
	defer func() { gooseUp = originalGooseUp }()

	os.Setenv("MIGRATIONS_DIR", "./fake_dir")

	err = migrate(mockDB)
	require.NoError(t, err)
}

func TestMigrateFailure(t *testing.T) {
	mockDB, _, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed initating test %v", err)
	}
	// Enable logging to track actual SQL queries
	goose.SetLogger(log.New(os.Stdout, "[GOOSE] ", log.LstdFlags))

	// Override gooseUp to simulate success
	originalGooseUp := gooseUp
	gooseUp = func(db *sql.DB, dir string, optionsFunc ...goose.OptionsFunc) error {
		return fmt.Errorf("fake migration failure")
	}
	defer func() { gooseUp = originalGooseUp }()

	os.Setenv("MIGRATIONS_DIR", "./fake_dir")

	err = migrate(mockDB)
	require.Error(t, err)
	require.Contains(t, err.Error(), "fake migration failure")
}
