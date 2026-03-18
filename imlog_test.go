package imlog

import (
	"os"
	"testing"
)

func TestLogger(t *testing.T) {
	tmpDir := t.TempDir()

	logger, err := NewLogger(tmpDir)
	if err != nil {
		t.Fatalf("Failed to create logger: %v", err)
	}
	defer logger.Close()

	tests := []struct {
		name    string
		message string
		fn      func(string) error
	}{
		{"Info", "This is an info message", logger.Info},
		{"Error", "This is an error message", logger.Error},
		{"Warning", "This is a warning message", logger.Warning},
		{"Event", "This is an event message", logger.Event},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.fn(tt.message)
			if err != nil {
				t.Errorf("Failed to log %s: %v", tt.name, err)
			}
		})
	}

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

	logger1, err := NewLogger(tmpDir)
	if err != nil {
		t.Fatalf("Failed to create logger: %v", err)
	}

	if err := logger1.Info("First message"); err != nil {
		t.Fatalf("Failed to log: %v", err)
	}

	hashAfterFirstLog := logger1.GetLastHashHex()
	logger1.Close()

	logger2, err := NewLogger(tmpDir)
	if err != nil {
		t.Fatalf("Failed to reopen logger: %v", err)
	}
	defer logger2.Close()

	if logger2.GetLastHashHex() != hashAfterFirstLog {
		t.Errorf("Hash was not preserved after reopening")
	}

	if err := logger2.Info("Second message"); err != nil {
		t.Fatalf("Failed to log second message: %v", err)
	}

	newHash := logger2.GetLastHashHex()
	if newHash == hashAfterFirstLog {
		t.Errorf("Hash should change after new log")
	}
}
