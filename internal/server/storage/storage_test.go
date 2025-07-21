package storage

import (
	"context"
	"testing"
	"time"

	"github.com/AlenaMolokova/gophkeeper/internal/shared/models"
	gomock "github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDefaultPoolConfig(t *testing.T) {
	config := DefaultPoolConfig()
	require.NotNil(t, config)

	assert.Equal(t, int32(20), config.MaxConns)
	assert.Equal(t, int32(5), config.MinConns)
	assert.Equal(t, 1*time.Hour, config.MaxConnLifetime)
	assert.Equal(t, 30*time.Minute, config.MaxConnIdleTime)
	assert.Equal(t, 1*time.Minute, config.HealthCheckPeriod)
}

func TestPoolStats(t *testing.T) {
	stats := &PoolStats{
		TotalConns:    10,
		IdleConns:     5,
		AcquiredConns: 3,
		Constructing:  2,
	}

	assert.Equal(t, int32(10), stats.TotalConns)
	assert.Equal(t, int32(5), stats.IdleConns)
	assert.Equal(t, int32(3), stats.AcquiredConns)
	assert.Equal(t, int32(2), stats.Constructing)
}

func TestNewPostgresStorageWithInvalidDSN(t *testing.T) {
	ctx := context.Background()
	config := DefaultPoolConfig()

	// Test with invalid DSN
	_, err := NewPostgresStorageWithConfig(ctx, "invalid-dsn", config)
	assert.Error(t, err)
}

func TestNewPostgresStorageWithNilConfig(t *testing.T) {
	ctx := context.Background()

	// Test with nil config
	_, err := NewPostgresStorageWithConfig(ctx, "postgres://test", nil)
	assert.Error(t, err)
}

func TestNewPostgresStorageWithInvalidConfig(t *testing.T) {
	ctx := context.Background()
	config := &PoolConfig{
		MaxConns: -1, // Invalid negative value
	}

	// Test with invalid config
	_, err := NewPostgresStorageWithConfig(ctx, "postgres://test", config)
	assert.Error(t, err)
}

// MockStorage tests using the existing mock.
func TestMockStorageInterface(t *testing.T) {
	// This test verifies that the mock implements the Storage interface correctly
	ctx := context.Background()

	// Create a mock storage instance with controller
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockStorage := NewMockStorage(ctrl)

	// Set up expectations for mock methods
	mockStorage.EXPECT().SaveUser(ctx, "test@example.com", "hash").Return("user-id", nil)
	mockStorage.EXPECT().FindUserByEmail(ctx, "test@example.com").Return(models.User{}, nil)
	mockStorage.EXPECT().SaveData(ctx, "user-id", models.Data{}).Return("data-id", nil)
	mockStorage.EXPECT().FindDataByID(ctx, "user-id", "data-id").Return(models.Data{}, nil)
	mockStorage.EXPECT().EditData(ctx, "user-id", models.Data{}).Return("data-id", nil)
	mockStorage.EXPECT().DeleteData(ctx, "user-id", "data-id").Return(nil)
	mockStorage.EXPECT().GetPoolStats().Return(&PoolStats{})

	// Test that the interface methods can be called
	_, err := mockStorage.SaveUser(ctx, "test@example.com", "hash")
	assert.NoError(t, err)

	_, err = mockStorage.FindUserByEmail(ctx, "test@example.com")
	assert.NoError(t, err)

	_, err = mockStorage.SaveData(ctx, "user-id", models.Data{})
	assert.NoError(t, err)

	_, err = mockStorage.FindDataByID(ctx, "user-id", "data-id")
	assert.NoError(t, err)

	_, err = mockStorage.EditData(ctx, "user-id", models.Data{})
	assert.NoError(t, err)

	err = mockStorage.DeleteData(ctx, "user-id", "data-id")
	assert.NoError(t, err)

	stats := mockStorage.GetPoolStats()
	assert.NotNil(t, stats)
}

func TestPoolConfigValidation(t *testing.T) {
	// Test various invalid configurations
	testCases := []struct {
		name   string
		config *PoolConfig
	}{
		{
			name: "negative MaxConns",
			config: &PoolConfig{
				MaxConns: -1,
				MinConns: 5,
			},
		},
		{
			name: "negative MinConns",
			config: &PoolConfig{
				MaxConns: 20,
				MinConns: -1,
			},
		},
		{
			name: "MinConns greater than MaxConns",
			config: &PoolConfig{
				MaxConns: 10,
				MinConns: 20,
			},
		},
		{
			name: "negative MaxConnLifetime",
			config: &PoolConfig{
				MaxConns:        20,
				MinConns:        5,
				MaxConnLifetime: -1 * time.Hour,
			},
		},
		{
			name: "negative MaxConnIdleTime",
			config: &PoolConfig{
				MaxConns:        20,
				MinConns:        5,
				MaxConnIdleTime: -1 * time.Minute,
			},
		},
		{
			name: "negative HealthCheckPeriod",
			config: &PoolConfig{
				MaxConns:          20,
				MinConns:          5,
				HealthCheckPeriod: -1 * time.Second,
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctx := context.Background()
			_, err := NewPostgresStorageWithConfig(ctx, "postgres://test", tc.config)
			assert.Error(t, err)
		})
	}
}
