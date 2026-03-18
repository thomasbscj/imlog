package imlog

import (
	"os"
	"testing"
)

func TestLogger(t *testing.T) {
	// Create temporary directory for testing
	tmpDir := t.TempDir()

	// Create a new logger
	logger, err := NewLogger(tmpDir)
	if err != nil {
		t.Fatalf("Failed to create logger: %v", err)
	}
	defer logger.Close()

	// Test logging different levels
	tests := []struct {
		name    string
		level   LogLevel
		message string
		fn      func(string) error
	}{
		{"Info", InfoLevel, "This is an info message", logger.Info},
		{"Error", ErrorLevel, "This is an error message", logger.Error},
		{"Warning", WarningLevel, "This is a warning message", logger.Warning},
		{"Event", EventLevel, "This is an event message", logger.Event},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.fn(tt.message)
			if err != nil {
				t.Errorf("Failed to log %s: %v", tt.name, err)
			}
		})
	}

	// Verify files were created
	logsPath := tmpDir + "/logs.log"
	sumsPath := tmpDir + "/sum.log"

	if _, err := os.Stat(logsPath); err != nil {
		t.Errorf("logs.log was not created: %v", err)
	}

	if _, err := os.Stat(sumsPath); err != nil {
		t.Errorf("sum.log was not created: %v", err)
	}
}

func TestLoggerReopen(t *testing.T) {
	tmpDir := t.TempDir()

	// Create logger and add entries
	logger1, err := NewLogger(tmpDir)
	if err != nil {
		t.Fatalf("Failed to create logger: %v", err)
	}

	logger1.GetLastHash()
	if err := logger1.Info("First message"); err != nil {
		t.Fatalf("Failed to log: %v", err)
	}

	hashAfterFirstLog := logger1.GetLastHash()
	logger1.Close()

	// Reopen logger and verify state
	logger2, err := NewLogger(tmpDir)
	if err != nil {
		t.Fatalf("Failed to reopen logger: %v", err)
	}
	defer logger2.Close()

	// The hash should be restored from file
	if logger2.GetLastHash() != hashAfterFirstLog {
		t.Errorf("Hash was not preserved after reopening")
	}

	// Add another message
	if err := logger2.Info("Second message"); err != nil {
		t.Fatalf("Failed to log second message: %v", err)
	}

	newHash := logger2.GetLastHash()
	if newHash == hashAfterFirstLog {
		t.Errorf("Hash should change after new log")
	}
}
