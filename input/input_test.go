package input

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestGetFileInput(t *testing.T) {
	// Create a temporary directory for test files
	tempDir, err := os.MkdirTemp("", "input_test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	tests := []struct {
		name           string
		fileContent    string
		expectedOutput string
		expectError    bool
	}{
		{
			name:           "simple single line text",
			fileContent:    "Hello, World!",
			expectedOutput: "Hello, World!",
			expectError:    false,
		},
		{
			name:           "multi-line text",
			fileContent:    "Line 1\nLine 2\nLine 3",
			expectedOutput: "Line 1Line 2Line 3",
			expectError:    false,
		},
		{
			name:           "empty file",
			fileContent:    "",
			expectedOutput: "",
			expectError:    false,
		},
		{
			name:           "text with special characters",
			fileContent:    "Special chars: !@#$%^&*()_+-={}[]|\\:;\"'<>?,./",
			expectedOutput: "Special chars: !@#$%^&*()_+-={}[]|\\:;\"'<>?,./",
			expectError:    false,
		},
		{
			name:           "text with unicode characters",
			fileContent:    "Unicode: 你好世界 🌍 café naïve résumé",
			expectedOutput: "Unicode: 你好世界 🌍 café naïve résumé",
			expectError:    false,
		},
		{
			name:           "text with tabs and spaces",
			fileContent:    "Text\twith\ttabs\nand    spaces",
			expectedOutput: "Text\twith\ttabsand    spaces",
			expectError:    false,
		},
		{
			name:           "large text content",
			fileContent:    strings.Repeat("Large content line.\n", 1000),
			expectedOutput: strings.Repeat("Large content line.", 1000),
			expectError:    false,
		},
		{
			name:           "text with only whitespace lines",
			fileContent:    "   \n\t\n   \n",
			expectedOutput: "   \t   ",
			expectError:    false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create test file
			testFile := filepath.Join(tempDir, tt.name+".txt")
			err := os.WriteFile(testFile, []byte(tt.fileContent), 0644)
			if err != nil {
				t.Fatalf("Failed to create test file: %v", err)
			}

			// Test the function
			result, err := GetFileInput(testFile)

			// Check error expectation
			if tt.expectError && err == nil {
				t.Errorf("Expected error but got none")
			}
			if !tt.expectError && err != nil {
				t.Errorf("Unexpected error: %v", err)
			}

			// Check result
			if result != tt.expectedOutput {
				t.Errorf("Expected %q, got %q", tt.expectedOutput, result)
			}
		})
	}
}

func TestGetFileInputErrors(t *testing.T) {
	tests := []struct {
		name        string
		filePath    string
		expectError bool
		errorSubstr string
	}{
		{
			name:        "non-existent file",
			filePath:    "/path/that/does/not/exist.txt",
			expectError: true,
			errorSubstr: "couldn't open file",
		},
		{
			name:        "directory instead of file",
			filePath:    ".",
			expectError: true,
			errorSubstr: "error reading file",
		},
		{
			name:        "empty path",
			filePath:    "",
			expectError: true,
			errorSubstr: "couldn't open file",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := GetFileInput(tt.filePath)

			if tt.expectError {
				if err == nil {
					t.Errorf("Expected error but got none")
				}
				if !strings.Contains(err.Error(), tt.errorSubstr) {
					t.Errorf("Expected error to contain %q, got %q", tt.errorSubstr, err.Error())
				}
				if result != "" {
					t.Errorf("Expected empty result on error, got %q", result)
				}
			} else {
				if err != nil {
					t.Errorf("Unexpected error: %v", err)
				}
			}
		})
	}
}

func TestGetFileInputPermissions(t *testing.T) {
	// Skip this test on Windows as file permissions work differently
	if os.Getenv("GOOS") == "windows" {
		t.Skip("Skipping permission test on Windows")
	}

	// Create a temporary directory for test files
	tempDir, err := os.MkdirTemp("", "input_perm_test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Create a file with no read permissions
	testFile := filepath.Join(tempDir, "no_read_perm.txt")
	err = os.WriteFile(testFile, []byte("test content"), 0644)
	if err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	// Remove read permissions
	err = os.Chmod(testFile, 0000)
	if err != nil {
		t.Fatalf("Failed to change file permissions: %v", err)
	}

	// Test the function
	result, err := GetFileInput(testFile)

	// Should get an error
	if err == nil {
		t.Errorf("Expected error for file with no read permissions")
	}
	if result != "" {
		t.Errorf("Expected empty result on error, got %q", result)
	}
	if !strings.Contains(err.Error(), "couldn't open file") {
		t.Errorf("Expected error to contain 'couldn't open file', got %q", err.Error())
	}
}
