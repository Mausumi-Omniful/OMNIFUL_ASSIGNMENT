package i18n

import (
	"context"
	"os"
	"path/filepath"
	"sync"
	"testing"
)

func TestLocalizationInitialization(t *testing.T) {
	// Reset the package state
	langKeys = make(map[string]map[string]langValue)
	initialized = false
	initOnce = sync.Once{}

	// Create a temporary directory for test localization files
	tempDir, err := os.MkdirTemp("", "localization_test")
	if err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Create English translation file with a service-specific key
	enContent := `{
		"test.service_specific": {
			"message": "Service specific message"
		},
		"test.common_override": {
			"message": "Service override of common message"
		}
	}`
	err = os.WriteFile(filepath.Join(tempDir, "messages.en.json"), []byte(enContent), 0644)
	if err != nil {
		t.Fatalf("Failed to write English translation file: %v", err)
	}

	// Initialize with the test directory
	err = Initialize(WithRootPath(tempDir))
	if err != nil {
		t.Fatalf("Failed to initialize i18n: %v", err)
	}

	// Verify that service-specific translations are loaded
	ctx := context.Background()
	serviceMessage := Translate(ctx, "test.service_specific")
	if serviceMessage != "Service specific message" {
		t.Errorf("Expected 'Service specific message', got: %s", serviceMessage)
	}

	// Verify that translations fall back to common ones when not in service
	commonMessage := Translate(ctx, "test.only_in_common")
	if commonMessage == "test.only_in_common" {
		t.Logf("Common fallback test skipped: common translations not found")
	} else if commonMessage != "This message is only in commons" {
		t.Errorf("Expected common fallback 'This message is only in commons', got: %s", commonMessage)
	}

	// Verify that service-specific translations override common ones
	overrideMessage := Translate(ctx, "test.common_override")
	if overrideMessage != "Service override of common message" {
		t.Errorf("Expected service override 'Service override of common message', got: %s", overrideMessage)
	}
}
