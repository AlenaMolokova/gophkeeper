package tui

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestNewLocalizer tests localizer creation.
func TestNewLocalizer(t *testing.T) {
	// Skip if locales directory doesn't exist
	if _, err := os.Stat("locales"); os.IsNotExist(err) {
		t.Skip("Skipping test because locales directory doesn't exist")
	}

	localizer, err := NewLocalizer()
	require.NoError(t, err)
	assert.NotNil(t, localizer)
	assert.NotNil(t, localizer.bundle)
	assert.NotNil(t, localizer.localizer)
}

// TestNewLocalizerWithMissingLocalesDir tests localizer creation with missing locales directory.
func TestNewLocalizerWithMissingLocalesDir(t *testing.T) {
	// Temporarily rename locales directory
	originalName := "locales"
	tempName := "locales_temp"

	// Check if locales directory exists
	if _, err := os.Stat(originalName); err == nil {
		err := os.Rename(originalName, tempName)
		require.NoError(t, err)
		defer func() {
			// Restore original directory
			if _, err := os.Stat(tempName); err == nil {
				err := os.Rename(tempName, originalName)
				if err != nil {
					t.Logf("Failed to restore original directory: %v", err)
				}
			}
		}()
	}

	_, err := NewLocalizer()
	assert.Error(t, err)
}

// TestTranslate tests translation functionality.
func TestTranslate(t *testing.T) {
	// Skip if locales directory doesn't exist
	if _, err := os.Stat("locales"); os.IsNotExist(err) {
		t.Skip("Skipping test because locales directory doesn't exist")
	}

	localizer, err := NewLocalizer()
	require.NoError(t, err)

	// Test existing translation
	result := localizer.Translate("registration")
	assert.NotEqual(t, "registration", result) // Should be translated
	assert.NotEmpty(t, result)

	// Test non-existing translation (should return key)
	result = localizer.Translate("non_existing_key")
	assert.Equal(t, "non_existing_key", result)
}

// TestTranslateWithParams tests translation with parameters.
func TestTranslateWithParams(t *testing.T) {
	// Skip if locales directory doesn't exist
	if _, err := os.Stat("locales"); os.IsNotExist(err) {
		t.Skip("Skipping test because locales directory doesn't exist")
	}

	localizer, err := NewLocalizer()
	require.NoError(t, err)

	params := map[string]interface{}{
		"name": "test",
	}

	// Test translation with parameters
	result := localizer.TranslateWithParams("non_existing_key", params)
	assert.Equal(t, "non_existing_key", result) // Should return key for non-existing translation
}

// TestTranslateWithEmptyParams tests translation with empty parameters.
func TestTranslateWithEmptyParams(t *testing.T) {
	// Skip if locales directory doesn't exist
	if _, err := os.Stat("locales"); os.IsNotExist(err) {
		t.Skip("Skipping test because locales directory doesn't exist")
	}

	localizer, err := NewLocalizer()
	require.NoError(t, err)

	params := map[string]interface{}{}

	result := localizer.TranslateWithParams("registration", params)
	assert.NotEqual(t, "registration", result)
	assert.NotEmpty(t, result)
}

// TestLocalizerLanguageDetection tests language detection from environment.
func TestLocalizerLanguageDetection(t *testing.T) {
	// Skip if locales directory doesn't exist
	if _, err := os.Stat("locales"); os.IsNotExist(err) {
		t.Skip("Skipping test because locales directory doesn't exist")
	}

	// Test Russian (default)
	os.Setenv("LANG", "ru_RU.UTF-8")
	localizer, err := NewLocalizer()
	require.NoError(t, err)
	assert.NotNil(t, localizer)

	// Test English
	os.Setenv("LANG", "en_US.UTF-8")
	localizer, err = NewLocalizer()
	require.NoError(t, err)
	assert.NotNil(t, localizer)

	// Test empty LANG (should default to Russian)
	os.Setenv("LANG", "")
	localizer, err = NewLocalizer()
	require.NoError(t, err)
	assert.NotNil(t, localizer)

	// Test short language code
	os.Setenv("LANG", "en")
	localizer, err = NewLocalizer()
	require.NoError(t, err)
	assert.NotNil(t, localizer)
}

// TestLocalizerBundleLoading tests bundle loading functionality.
func TestLocalizerBundleLoading(t *testing.T) {
	// Skip if locales directory doesn't exist
	if _, err := os.Stat("locales"); os.IsNotExist(err) {
		t.Skip("Skipping test because locales directory doesn't exist")
	}

	localizer, err := NewLocalizer()
	require.NoError(t, err)

	// Test that bundle is properly loaded
	assert.NotNil(t, localizer.bundle)

	// Test that localizer is properly configured
	assert.NotNil(t, localizer.localizer)
}

// TestLocalizerMultipleTranslations tests multiple translation calls.
func TestLocalizerMultipleTranslations(t *testing.T) {
	// Skip if locales directory doesn't exist
	if _, err := os.Stat("locales"); os.IsNotExist(err) {
		t.Skip("Skipping test because locales directory doesn't exist")
	}

	localizer, err := NewLocalizer()
	require.NoError(t, err)

	// Test multiple translations
	translations := []string{
		"registration",
		"login",
		"view_data",
		"add_data",
		"edit_data",
		"delete_data",
		"exit",
	}

	for _, key := range translations {
		result := localizer.Translate(key)
		assert.NotEqual(t, key, result) // Should be translated
		assert.NotEmpty(t, result)
	}
}

// TestLocalizerErrorHandling tests error handling in translation.
func TestLocalizerErrorHandling(t *testing.T) {
	// Skip if locales directory doesn't exist
	if _, err := os.Stat("locales"); os.IsNotExist(err) {
		t.Skip("Skipping test because locales directory doesn't exist")
	}

	localizer, err := NewLocalizer()
	require.NoError(t, err)

	// Test with empty key
	result := localizer.Translate("")
	assert.Equal(t, "", result)

	// Test with special characters
	result = localizer.Translate("key_with_special_chars_!@#$%")
	assert.Equal(t, "key_with_special_chars_!@#$%", result)
}

// TestLocalizerWithNilBundle tests localizer behavior with nil bundle.
func TestLocalizerWithNilBundle(t *testing.T) {
	localizer := &Localizer{
		bundle: nil,
	}

	// Test that Translate doesn't panic with nil bundle
	result := localizer.Translate("test.key")
	assert.Equal(t, "test.key", result) // Should return the key as fallback

	// Test that TranslateWithParams doesn't panic with nil bundle
	result = localizer.TranslateWithParams("test.key", map[string]interface{}{"param": "value"})
	assert.Equal(t, "test.key", result) // Should return the key as fallback
}

// TestLocalizerWithEmptyKey tests localizer behavior with empty keys.
func TestLocalizerWithEmptyKey(t *testing.T) {
	localizer := &Localizer{
		bundle: nil,
	}

	// Test with empty key
	result := localizer.Translate("")
	assert.Equal(t, "", result)

	// Test with empty key and params
	result = localizer.TranslateWithParams("", map[string]interface{}{"param": "value"})
	assert.Equal(t, "", result)
}

// TestLocalizerWithNilParams tests localizer behavior with nil parameters.
func TestLocalizerWithNilParams(t *testing.T) {
	localizer := &Localizer{
		bundle: nil,
	}

	// Test with nil params
	result := localizer.TranslateWithParams("test.key", nil)
	assert.Equal(t, "test.key", result)
}

// TestLocalizerWithEmptyParams tests localizer behavior with empty parameters.
func TestLocalizerWithEmptyParams(t *testing.T) {
	localizer := &Localizer{
		bundle: nil,
	}

	// Test with empty params map
	result := localizer.TranslateWithParams("test.key", map[string]interface{}{})
	assert.Equal(t, "test.key", result)
}

// TestLocalizerConcurrentAccess tests localizer behavior under concurrent access.
func TestLocalizerConcurrentAccess(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping concurrent test in short mode")
	}

	localizer := &Localizer{
		bundle: nil,
	}

	// Test concurrent access to Translate
	done := make(chan bool, 10)
	for i := 0; i < 10; i++ {
		go func() {
			defer func() { done <- true }()
			result := localizer.Translate("test.key")
			assert.Equal(t, "test.key", result)
		}()
	}

	// Wait for all goroutines to complete
	for i := 0; i < 10; i++ {
		<-done
	}

	// Test concurrent access to TranslateWithParams
	for i := 0; i < 10; i++ {
		go func() {
			defer func() { done <- true }()
			result := localizer.TranslateWithParams("test.key", map[string]interface{}{"param": "value"})
			assert.Equal(t, "test.key", result)
		}()
	}

	// Wait for all goroutines to complete
	for i := 0; i < 10; i++ {
		<-done
	}
}

// TestLocalizerTranslationConsistency tests translation consistency.
func TestLocalizerTranslationConsistency(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping translation consistency test in short mode")
	}

	localizer := &Localizer{
		bundle: nil,
	}

	// Test that same key returns same result
	key := "test.key"
	result1 := localizer.Translate(key)
	result2 := localizer.Translate(key)
	assert.Equal(t, result1, result2)

	// Test that same key with same params returns same result
	params := map[string]interface{}{"param": "value"}
	result3 := localizer.TranslateWithParams(key, params)
	result4 := localizer.TranslateWithParams(key, params)
	assert.Equal(t, result3, result4)
}

// TestLocalizerFallbackBehavior tests localizer fallback behavior.
func TestLocalizerFallbackBehavior(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping fallback behavior test in short mode")
	}

	localizer := &Localizer{
		bundle: nil,
	}

	// Test fallback behavior when translation is not found
	result := localizer.Translate("nonexistent.key")
	assert.Equal(t, "nonexistent.key", result) // Should return the key as fallback

	// Test fallback behavior with params
	result = localizer.TranslateWithParams("nonexistent.key", map[string]interface{}{"param": "value"})
	assert.Equal(t, "nonexistent.key", result) // Should return the key as fallback
}
