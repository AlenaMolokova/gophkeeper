package sync

import (
	"encoding/json"
	"testing"
	"time"

	clientapi "github.com/AlenaMolokova/gophkeeper/pkg/client/api"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewDefaultSyncer(t *testing.T) {
	syncer, err := NewDefaultSyncer()
	require.NoError(t, err)
	require.NotNil(t, syncer)
	require.NotNil(t, syncer.Storage)
	syncer.Storage.Close()
}

func TestNewDefaultSyncerWithPath(t *testing.T) {
	tempDir := t.TempDir()
	dbPath := tempDir + "/test.db"

	syncer, err := NewDefaultSyncerWithPath(dbPath)
	require.NoError(t, err)
	require.NotNil(t, syncer)
	require.NotNil(t, syncer.Storage)
	syncer.Storage.Close()
}

func TestMergeFunction(t *testing.T) {
	// Create syncer with unique path to avoid file locking issues
	tempDir := t.TempDir()
	dbPath := tempDir + "/merge-test.db"
	syncer, err := NewDefaultSyncerWithPath(dbPath)
	require.NoError(t, err)
	defer syncer.Storage.Close()

	// Create test data with timestamps
	now := time.Now().Unix()
	localData1 := clientapi.Data{Id: "id1", Type: "login", Payload: []byte("local-data-1"), Timestamp: now}
	localData2 := clientapi.Data{Id: "id2", Type: "login", Payload: []byte("local-data-2"), Timestamp: now}
	remoteData2 := clientapi.Data{Id: "id2", Type: "login", Payload: []byte("remote-data-2"), Timestamp: now + 1} // Newer
	remoteData3 := clientapi.Data{Id: "id3", Type: "login", Payload: []byte("remote-data-3"), Timestamp: now}

	// Marshal to JSON
	local1Bytes, _ := json.Marshal(localData1)
	local2Bytes, _ := json.Marshal(localData2)
	remote2Bytes, _ := json.Marshal(remoteData2)
	remote3Bytes, _ := json.Marshal(remoteData3)

	local := map[string][]byte{
		"id1": local1Bytes,
		"id2": local2Bytes,
	}

	remote := map[string][]byte{
		"id2": remote2Bytes, // Should override local due to newer timestamp
		"id3": remote3Bytes,
	}

	merged := syncer.Merge(local, remote)

	// Should contain all unique keys
	assert.Len(t, merged, 3)

	// Check id1 (only in local)
	var data1 clientapi.Data
	err = json.Unmarshal(merged["id1"], &data1)
	require.NoError(t, err)
	assert.Equal(t, "local-data-1", string(data1.Payload))

	// Check id2 (remote should override local due to newer timestamp)
	var data2 clientapi.Data
	err = json.Unmarshal(merged["id2"], &data2)
	require.NoError(t, err)
	assert.Equal(t, "remote-data-2", string(data2.Payload))

	// Check id3 (only in remote)
	var data3 clientapi.Data
	err = json.Unmarshal(merged["id3"], &data3)
	require.NoError(t, err)
	assert.Equal(t, "remote-data-3", string(data3.Payload))
}

func TestMergeWithEmptyMaps(t *testing.T) {
	// Create syncer with unique path to avoid file locking issues
	tempDir := t.TempDir()
	dbPath := tempDir + "/empty-test.db"
	syncer, err := NewDefaultSyncerWithPath(dbPath)
	require.NoError(t, err)
	defer syncer.Storage.Close()

	// Test with empty local
	merged := syncer.Merge(map[string][]byte{}, map[string][]byte{"id1": []byte("data1")})
	assert.Len(t, merged, 1)
	assert.Equal(t, []byte("data1"), merged["id1"])

	// Test with empty remote
	merged = syncer.Merge(map[string][]byte{"id1": []byte("data1")}, map[string][]byte{})
	assert.Len(t, merged, 1)
	assert.Equal(t, []byte("data1"), merged["id1"])

	// Test with both empty
	merged = syncer.Merge(map[string][]byte{}, map[string][]byte{})
	assert.Len(t, merged, 0)
}
