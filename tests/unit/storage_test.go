package unit

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/AlenaMolokova/gophkeeper/internal/client/storage"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetDefaultDBPath(t *testing.T) {
	path := storage.GetDefaultDBPath()
	assert.NotEmpty(t, path)
	assert.Contains(t, path, ".gophkeeper")
	assert.Contains(t, path, "data.db")
}

func TestEnsureDBDirectory(t *testing.T) {
	tempDir := t.TempDir()
	dbPath := filepath.Join(tempDir, "test", "data.db")

	err := storage.EnsureDBDirectory(dbPath)
	require.NoError(t, err)

	// Check that directory was created
	dir := filepath.Dir(dbPath)
	_, err = os.Stat(dir)
	require.NoError(t, err)
}

func TestNewStorage(t *testing.T) {
	tempDir := t.TempDir()
	dbPath := filepath.Join(tempDir, "test.db")

	st, err := storage.NewStorage(dbPath)
	require.NoError(t, err)
	require.NotNil(t, st)
	defer st.Close()

	// Check that database file was created
	_, err = os.Stat(dbPath)
	require.NoError(t, err)
}

func TestStorageOperations(t *testing.T) {
	tempDir := t.TempDir()
	dbPath := filepath.Join(tempDir, "test.db")

	st, err := storage.NewStorage(dbPath)
	require.NoError(t, err)
	defer st.Close()

	// Test SaveData
	testData := []byte("test data")
	err = st.SaveData("test-id", testData)
	require.NoError(t, err)

	// Test GetData
	retrievedData, err := st.GetData("test-id")
	require.NoError(t, err)
	assert.Equal(t, testData, retrievedData)

	// Test GetData with non-existent ID
	_, err = st.GetData("non-existent")
	assert.Error(t, err)

	// Test GetAllData
	allData, err := st.GetAllData()
	require.NoError(t, err)
	assert.Len(t, allData, 1)
	assert.Equal(t, testData, allData["test-id"])
}

func TestStorageClose(t *testing.T) {
	tempDir := t.TempDir()
	dbPath := filepath.Join(tempDir, "test.db")

	st, err := storage.NewStorage(dbPath)
	require.NoError(t, err)

	// Test that Close doesn't panic
	assert.NotPanics(t, func() {
		st.Close()
	})
}
