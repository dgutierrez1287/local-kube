package settings

import (
	"os"
  "os/exec"
	"path/filepath"
	"testing"

  "github.com/dgutierrez1287/local-kube/logger"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)


// TestMain is executed before running any tests
func TestMain(m *testing.M) {
	// Initialize the logger before running any tests
	logger.InitLogging(false, true, false)
	os.Exit(m.Run())
}

// Test GetAppDirPath
func TestGetAppDirPath(t *testing.T) {
	homeDir, err := os.UserHomeDir()
	require.NoError(t, err)

	expectedPath := filepath.Join(homeDir, ".local-kube")
	assert.Equal(t, expectedPath, GetAppDirPath())
}

// Test DirectoryExists when directory exists
func TestDirectoryExists_Exists(t *testing.T) {
	// Create temp directory
	tempDir := t.TempDir()

	assert.True(t, DirectoryExists(tempDir))
}

// Test DirectoryExists when directory does not exist
func TestDirectoryExists_NotExists(t *testing.T) {
	nonExistentDir := filepath.Join(os.TempDir(), "does-not-exist-12345")

	assert.False(t, DirectoryExists(nonExistentDir))
}

// Test CreateDirectory successfully
func TestCreateDirectory_Success(t *testing.T) {
	tempDir := filepath.Join(os.TempDir(), "test-create-dir")
	defer os.Remove(tempDir) // Cleanup

	CreateDirectory(tempDir)

	assert.DirExists(t, tempDir)
}

// TestCreateDirectory_Failure tests that CreateDirectory calls os.Exit(123) when it fails
func TestCreateDirectory_Failure(t *testing.T) {
	if os.Geteuid() == 0 {
		t.Skip("Skipping test since it's running as root.")
	}

	// Run test in a subprocess
	cmd := exec.Command(os.Args[0], "-test.run=TestHelperProcess")
	cmd.Env = append(os.Environ(), "TEST_CRASH=1")
	err := cmd.Run()

	// Assert that os.Exit(123) was triggered
	exitError, ok := err.(*exec.ExitError)
	assert.True(t, ok, "Process should have exited with an error")
	assert.Equal(t, 123, exitError.ExitCode(), "Expected os.Exit(123) to be called on failure")
}

// TestHelperProcess is a helper function to trigger os.Exit in a subprocess
func TestHelperProcess(t *testing.T) {
	if os.Getenv("TEST_CRASH") == "1" {
		CreateDirectory("/root/forbidden-dir") // Should trigger os.Exit(123)
	}
}
