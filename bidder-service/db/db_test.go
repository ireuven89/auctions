package db

/*func TestMigrate(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	// goose typically starts a transaction
	mock.ExpectBegin()

	// goose checks current migration version
	mock.ExpectQuery("(?i)^SELECT version_id, is_applied FROM goose_db_version ORDER BY id DESC$").
		WillReturnRows(sqlmock.NewRows([]string{"version_id", "is_applied"}))

	// goose may apply migrations (INSERT version, etc.)
	mock.ExpectExec("(?i)^INSERT INTO goose_db_version").
		WillReturnResult(sqlmock.NewResult(1, 1))

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
*/
