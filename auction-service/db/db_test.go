package db

//
//import (
//	"database/sql"
//	"database/sql/driver"
//	"fmt"
//	"github.com/DATA-DOG/go-sqlmock"
//	"github.com/sethvargo/go-retry"
//	"github.com/stretchr/testify/assert"
//	"github.com/stretchr/testify/require"
//	"testing"
//	"time"
//)
//
//var mockDB = &sql.DB{}
//
//func setup() {
//	mockDb, _, err := sqlmock.New()
//
//	if err != nil {
//		panic(err)
//	}
//
//	mockDB = mockDb
//}
//
//type mockDriver struct {
//	openFunc func(name string) (driver.Conn, error)
//}
//
//func (d *mockDriver) Open(name string) (driver.Conn, error) {
//	if d.openFunc != nil {
//		return d.openFunc(name)
//	}
//	return nil, fmt.Errorf("mock connection error")
//}
//
//func TestConnectWithRetry_Success(t *testing.T) {
//	// Create a mock database
//	db, mock, err := sqlmock.New()
//	require.NoError(t, err)
//	defer db.Close()
//
//	// Set expectations for successful flow
//	mock.ExpectPing()
//	mock.ExpectExec("CREATE DATABASE IF NOT EXISTS `testdb`").
//		WillReturnResult(sqlmock.NewResult(0, 0))
//	mock.ExpectPing() // Second ping for database-specific connection
//
//	// Override sql.Open for this test
//	originalOpen := sql.Open
//	callCount := 0
//	sql.Open = func(driverName, dataSourceName string) (*sql.DB, error) {
//		callCount++
//		return db, nil
//	}
//	defer func() { sql.Open = originalOpen }()
//
//	result, err := MustNewDB("user", "pass", "localhost", 3306)
//
//	assert.NoError(t, err)
//	assert.NotNil(t, result)
//	assert.Equal(t, 2, callCount) // Should be called twice (initial + with database)
//	assert.NoError(t, mock.ExpectationsWereMet())
//}
//
//func TestConnectWithRetry_PingFailureWithRetry(t *testing.T) {
//	db, mock, err := sqlmock.New()
//	require.NoError(t, err)
//	defer db.Close()
//
//	// First two pings fail, third succeeds
//	mock.ExpectPing().WillReturnError(fmt.Errorf("connection refused"))
//	mock.ExpectPing().WillReturnError(fmt.Errorf("connection refused"))
//	mock.ExpectPing()
//	mock.ExpectExec("CREATE DATABASE IF NOT EXISTS `testdb`").
//		WillReturnResult(sqlmock.NewResult(0, 0))
//	mock.ExpectPing() // Final ping for database-specific connection
//
//	originalOpen := sql.Open
//	sql.Open = func(driverName, dataSourceName string) (*sql.DB, error) {
//		return db, nil
//	}
//	defer func() { sql.Open = originalOpen }()
//
//	result, err := MustNewDB("user", "pass", "localhost", 3306)
//
//	assert.NoError(t, err)
//	assert.NotNil(t, result)
//	assert.NoError(t, mock.ExpectationsWereMet())
//}
//
//func TestConnectWithRetry_DatabaseCreationFailure(t *testing.T) {
//	db, mock, err := sqlmock.New()
//	require.NoError(t, err)
//	defer db.Close()
//
//	mock.ExpectPing()
//	mock.ExpectExec("CREATE DATABASE IF NOT EXISTS `testdb`").
//		WillReturnError(fmt.Errorf("access denied"))
//
//	originalOpen := sql.Open
//	sql.Open = func(driverName, dataSourceName string) (*sql.DB, error) {
//		return db, nil
//	}
//	defer func() { sql.Open = originalOpen }()
//
//	result, err := connectWithRetry("user", "pass", "localhost", 3306, "testdb")
//
//	assert.Error(t, err)
//	assert.Nil(t, result)
//	assert.Contains(t, err.Error(), "failed creating database")
//	assert.NoError(t, mock.ExpectationsWereMet())
//}
//
//func TestConnectWithRetry_MaxRetriesExceeded(t *testing.T) {
//	db, mock, err := sqlmock.New()
//	require.NoError(t, err)
//	defer db.Close()
//
//	// Always fail ping to exceed max retries
//	for i := 0; i < 6; i++ { // More than max retries (5)
//		mock.ExpectPing().WillReturnError(fmt.Errorf("connection refused"))
//	}
//
//	originalOpen := sql.Open
//	sql.Open = func(driverName, dataSourceName string) (*sql.DB, error) {
//		return db, nil
//	}
//	defer func() { sql.Open = originalOpen }()
//
//	result, err := connectWithRetry("user", "pass", "localhost", 3306, "testdb")
//
//	assert.Error(t, err)
//	assert.Nil(t, result)
//	assert.Contains(t, err.Error(), "failed pinging MySQL")
//}
//
//func TestConnectWithRetry_SQLOpenFailure(t *testing.T) {
//	originalOpen := sql.Open
//	attemptCount := 0
//	sql.Open = func(driverName, dataSourceName string) (*sql.DB, error) {
//		attemptCount++
//		if attemptCount <= 2 {
//			return nil, fmt.Errorf("driver not found")
//		}
//		// Success on third attempt
//		db, mock, _ := sqlmock.New()
//		mock.ExpectPing()
//		mock.ExpectExec("CREATE DATABASE IF NOT EXISTS `testdb`").
//			WillReturnResult(sqlmock.NewResult(0, 0))
//		mock.ExpectPing()
//		return db, nil
//	}
//	defer func() { sql.Open = originalOpen }()
//
//	result, err := connectWithRetry("user", "pass", "localhost", 3306, "testdb")
//
//	assert.NoError(t, err)
//	assert.NotNil(t, result)
//	assert.Equal(t, 4, attemptCount) // Initial failures + 2 successful calls
//}
//
//func TestConnectWithRetry_ContextCancellation(t *testing.T) {
//	// Test with a very short context timeout
//	originalConnectFunc := connectWithRetry
//
//	// Create a version that accepts context
//	connectWithRetryCtx := func(ctx context.Context, user, password, host string, port int, databaseName string) (*sql.DB, error) {
//		dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/", user, password, host, port)
//		backoff := retry.WithMaxRetries(5, retry.NewExponential(3*time.Second))
//
//		return retry.DoValue[*sql.DB](ctx, backoff, func(ctx context.Context) (*sql.DB, error) {
//			return nil, retry.RetryableError(fmt.Errorf("always fail for test"))
//		})
//	}
//
//	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
//	defer cancel()
//
//	result, err := connectWithRetryCtx(ctx, "user", "pass", "localhost", 3306, "testdb")
//
//	assert.Error(t, err)
//	assert.Nil(t, result)
//	// Should fail due to context timeout, not max retries
//}
//
//func runner(m *testing.M) {
//	setup()
//}
/*
func TestMustNewDB(t *testing.T) {
	MustNewDB("failed db", "test", "test", 3306)
}
*/
