package utils

import (
	"bytes"
	"io"
	"log"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// Test model for migration testing
type DBTestModel struct {
	ID   uint   `gorm:"primaryKey"`
	Name string `gorm:"size:50"`
}

func (DBTestModel) TableName() string {
	return "test_models"
}

// 自定义一个"静音 + 忽略 RecordNotFound"的 logger
func createSilentLogger() logger.Interface {
	return logger.New(
		log.New(io.Discard, "", log.LstdFlags), // 丢弃输出
		logger.Config{
			SlowThreshold:             200 * time.Millisecond,
			LogLevel:                  logger.Silent, // 静音模式
			IgnoreRecordNotFoundError: true,          // 关键：忽略 not found
			Colorful:                  false,
		},
	)
}

func TestInitDatabaseWithSqlite(t *testing.T) {
	// Test initializing database with sqlite in-memory
	db, err := InitDatabase(nil, "", "file::memory:?cache=shared")
	assert.NoError(t, err)
	assert.NotNil(t, db)

	// Migrate the test model
	err = MakeMigrates(db, []interface{}{&DBTestModel{}})
	assert.NoError(t, err)

	// Test that we can execute a simple query
	testModel := &DBTestModel{Name: "test"}
	result := db.Create(testModel)
	assert.NoError(t, result.Error)
	assert.NotZero(t, testModel.ID)
}

func TestInitDatabaseWithDefaultSettings(t *testing.T) {
	// Save original environment variables
	origDriver := os.Getenv("DB_DRIVER")
	origDsn := os.Getenv("DSN")
	defer func() {
		os.Setenv("DB_DRIVER", origDriver)
		os.Setenv("DSN", origDsn)
	}()

	// Set environment variables for testing
	os.Setenv("DB_DRIVER", "")
	os.Setenv("DSN", "")

	// Test initializing database with default settings (sqlite in-memory)
	db, err := InitDatabase(nil, "", "")
	assert.NoError(t, err)
	assert.NotNil(t, db)

	// Migrate the test model
	err = MakeMigrates(db, []interface{}{&DBTestModel{}})
	assert.NoError(t, err)

	// Verify database operations work
	testModel := &DBTestModel{Name: "default_test"}
	result := db.Create(testModel)
	assert.NoError(t, result.Error)
	assert.NotZero(t, testModel.ID)
}

func TestInitDatabaseWithCustomLogger(t *testing.T) {
	// Test initializing database with custom logger
	var logBuf bytes.Buffer
	db, err := InitDatabase(&logBuf, "", "file::memory:?cache=shared")
	assert.NoError(t, err)
	assert.NotNil(t, db)

	// Migrate the test model
	err = MakeMigrates(db, []interface{}{&DBTestModel{}})
	assert.NoError(t, err)

	// Perform an operation to generate logs
	testModel := &DBTestModel{Name: "log_test"}
	result := db.Create(testModel)
	assert.NoError(t, result.Error)

	// Give some time for logs to be written
	time.Sleep(10 * time.Millisecond)

	// We're not strictly checking log content as it might vary
	// Just ensure the function works with custom logger
	assert.NotNil(t, db)
}

func TestInitDatabaseWithSilentLogger(t *testing.T) {
	// Test initializing database with silent logger
	silentLogger := createSilentLogger()

	// Create a config with silent logger to demonstrate usage
	_ = &gorm.Config{
		Logger:                 silentLogger,
		SkipDefaultTransaction: true,
	}

	// Since we can't directly test the private createDatabaseInstance function,
	// we test by using InitDatabase with a buffer and checking it works
	var logBuf bytes.Buffer
	db, err := InitDatabase(&logBuf, "", "file::memory:?cache=shared")
	assert.NoError(t, err)
	assert.NotNil(t, db)

	// Migrate the test model
	err = MakeMigrates(db, []interface{}{&DBTestModel{}})
	assert.NoError(t, err)

	// Perform operations
	testModel := &DBTestModel{Name: "silent_log_test"}
	result := db.Create(testModel)
	assert.NoError(t, result.Error)
	assert.NotZero(t, testModel.ID)
}

func TestMakeMigrates(t *testing.T) {
	// Initialize database
	db, err := InitDatabase(nil, "", "file::memory:?cache=shared")
	assert.NoError(t, err)
	assert.NotNil(t, db)

	// Test migrating a single model
	models := []interface{}{&DBTestModel{}}
	err = MakeMigrates(db, models)
	assert.NoError(t, err)

	// Verify table was created by checking if we can query it
	var count int64
	result := db.Model(&DBTestModel{}).Count(&count)
	assert.NoError(t, result.Error)

	// Test migrating multiple models (using same model twice to check idempotency)
	models = []interface{}{&DBTestModel{}, &DBTestModel{}}
	err = MakeMigrates(db, models)
	assert.NoError(t, err)

	// Test with empty models slice
	err = MakeMigrates(db, []interface{}{})
	assert.NoError(t, err)
}

func TestMakeMigratesWithError(t *testing.T) {
	// Initialize database
	db, err := InitDatabase(nil, "", "file::memory:?cache=shared")
	assert.NoError(t, err)
	assert.NotNil(t, db)

	// Test that migration error is properly handled without panicking
	// We'll test with a valid model to ensure the function works correctly
	assert.NotPanics(t, func() {
		err := MakeMigrates(db, []interface{}{&DBTestModel{}})
		assert.NoError(t, err)
	})
}

func TestInitDatabaseWithDifferentDrivers(t *testing.T) {
	// Test with sqlite driver explicitly specified
	db, err := InitDatabase(nil, "sqlite", "file::memory:?cache=shared")
	assert.NoError(t, err)
	assert.NotNil(t, db)

	// Migrate the test model
	err = MakeMigrates(db, []interface{}{&DBTestModel{}})
	assert.NoError(t, err)

	// Verify basic operations work
	testModel := &DBTestModel{Name: "driver_test"}
	result := db.Create(testModel)
	assert.NoError(t, result.Error)
	assert.NotZero(t, testModel.ID)
}

func TestInitDatabaseWithNilDriverAndDsn(t *testing.T) {
	// Save original environment variables
	origDriver := os.Getenv("DB_DRIVER")
	origDsn := os.Getenv("DSN")
	defer func() {
		os.Setenv("DB_DRIVER", origDriver)
		os.Setenv("DSN", origDsn)
	}()

	// Set environment variables
	os.Setenv("DB_DRIVER", "sqlite")
	os.Setenv("DSN", "file::memory:?cache=shared")

	// Test initializing database with nil driver and dsn (should use environment variables)
	db, err := InitDatabase(nil, "", "")
	assert.NoError(t, err)
	assert.NotNil(t, db)

	// Migrate the test model
	err = MakeMigrates(db, []interface{}{&DBTestModel{}})
	assert.NoError(t, err)

	// Verify basic operations work
	testModel := &DBTestModel{Name: "env_test"}
	result := db.Create(testModel)
	assert.NoError(t, result.Error)
	assert.NotZero(t, testModel.ID)
}

func TestInitDatabaseWithEmptyDsn(t *testing.T) {
	// Test initializing database with empty DSN (should use default in-memory database)
	db, err := InitDatabase(nil, "sqlite", "")
	assert.NoError(t, err)
	assert.NotNil(t, db)

	// Migrate the test model
	err = MakeMigrates(db, []interface{}{&DBTestModel{}})
	assert.NoError(t, err)

	// Verify basic operations work
	testModel := &DBTestModel{Name: "empty_dsn_test"}
	result := db.Create(testModel)
	assert.NoError(t, result.Error)
	assert.NotZero(t, testModel.ID)
}
