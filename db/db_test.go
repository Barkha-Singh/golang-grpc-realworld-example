package db

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDSN_Success(t *testing.T) {
	// Arrange: set all required env vars
	os.Setenv("DB_HOST", "myhost")
	os.Setenv("DB_USER", "myuser")
	os.Setenv("DB_PASSWORD", "mypassword")
	os.Setenv("DB_NAME", "mydb")
	os.Setenv("DB_PORT", "1234")
	defer func() {
		// clean up
		os.Unsetenv("DB_HOST")
		os.Unsetenv("DB_USER")
		os.Unsetenv("DB_PASSWORD")
		os.Unsetenv("DB_NAME")
		os.Unsetenv("DB_PORT")
	}()

	// Act
	actual, err := dsn()

	// Assert
	assert.NoError(t, err)
	expectedOptions := "charset=utf8mb4&parseTime=True&loc=Local"
	expected := "myuser:mypassword@(myhost:1234)/mydb?" + expectedOptions
	assert.Equal(t, expected, actual)
}

func TestDSN_MissingEnv(t *testing.T) {
	// Make sure DB_HOST is not set
	os.Unsetenv("DB_HOST")

	_, err := dsn()
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "$DB_HOST is not set")
}
