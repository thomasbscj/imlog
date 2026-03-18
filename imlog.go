// Immutable Logger - provides tamper-evident logging using SHA256 hashing
package imlog

import (
	"crypto/sha256"
	"encoding/binary"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"
)

type LogLevel string

const (
	InfoLevel    LogLevel = "info"
	ErrorLevel   LogLevel = "error"
	WarningLevel LogLevel = "warning"
	EventLevel   LogLevel = "event"
)


type LogEntry struct {
	Type      LogLevel  `json:"type"`
	Message   string    `json:"message"`
	Timestamp time.Time `json:"timestamp"`
	SHA256    [32]byte  `json:"sha256"`
}

type ImLogger struct {
	logsFile  *os.File
	sumsFile  *os.File
	logsPath  string
	sumsPath  string
	lastSum   [32]byte
	lastSumOffset int64
}

func NewLogger(dir string) (*ImLogger, error) {
	if err := os.MkdirAll(dir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create directory: %w", err)
	}

	logsPath := filepath.Join(dir, "logs.log")
	sumsPath := filepath.Join(dir, "sum.log")


	sumsFile, err := os.OpenFile(sumsPath, os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		return nil, fmt.Errorf("failed to open sum.log: %w", err)
	}

	logsFile, err := os.OpenFile(logsPath, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		sumsFile.Close()
		return nil, fmt.Errorf("failed to open logs.log: %w", err)
	}

	logger := &ImLogger{
		logsFile: logsFile,
		sumsFile: sumsFile,
		logsPath: logsPath,
		sumsPath: sumsPath,
	}

	if err := logger.initializeLastSum(); err != nil {
		logger.Close()
		return nil, err
	}

	return logger, nil
}

func (l *ImLogger) initializeLastSum() error {
	info, err := l.sumsFile.Stat()
	if err != nil {
		return fmt.Errorf("failed to stat sum.log: %w", err)
	}

	if info.Size() == 0 {
		offsetBuf := make([]byte, 8)
		binary.BigEndian.PutUint64(offsetBuf, 8) 
		if _, err := l.sumsFile.Write(offsetBuf); err != nil {
			return fmt.Errorf("failed to write initial offset: %w", err)
		}

		l.lastSum = [32]byte{}
		l.lastSumOffset = 8
		return nil
	}

	offsetBuf := make([]byte, 8)
	if _, err := l.sumsFile.ReadAt(offsetBuf, 0); err != nil {
		return fmt.Errorf("failed to read offset: %w", err)
	}

	offset := int64(binary.BigEndian.Uint64(offsetBuf))
	l.lastSumOffset = offset

	hashBuf := make([]byte, 32)
	n, err := l.sumsFile.ReadAt(hashBuf, offset)
	if err != nil && err != io.EOF {
		return fmt.Errorf("failed to read last sum: %w", err)
	}

	if n == 32 {
		copy(l.lastSum[:], hashBuf)
	}

	return nil
}

func (l *ImLogger) writeHash(hash [32]byte) error {
	newOffset := l.lastSumOffset + 32
	if _, err := l.sumsFile.WriteAt(hash[:], l.lastSumOffset); err != nil {
		return fmt.Errorf("failed to write hash: %w", err)
	}

	offsetBuf := make([]byte, 8)
	binary.BigEndian.PutUint64(offsetBuf, uint64(newOffset))
	if _, err := l.sumsFile.WriteAt(offsetBuf, 0); err != nil {
		return fmt.Errorf("failed to update offset: %w", err)
	}

	l.lastSum = hash
	l.lastSumOffset = newOffset

	return nil
}

func (l *ImLogger) Log(level LogLevel, message string) error {
	if l.logsFile == nil || l.sumsFile == nil {
		return errors.New("logger is closed")
	}

	entry := LogEntry{
		Type:      level,
		Message:   message,
		Timestamp: time.Now(),
	}

	entryJSON, err := json.Marshal(map[string]interface{}{
		"type":      entry.Type,
		"message":   entry.Message,
		"timestamp": entry.Timestamp,
	})
	if err != nil {
		return fmt.Errorf("failed to marshal log entry: %w", err)
	}

	data := append(entryJSON, l.lastSum[:]...)
	entry.SHA256 = sha256.Sum256(data)

	completeJSON, err := json.Marshal(entry)
	if err != nil {
		return fmt.Errorf("failed to marshal complete entry: %w", err)
	}

	if _, err := l.logsFile.Write(completeJSON); err != nil {
		return fmt.Errorf("failed to write log entry: %w", err)
	}

	if _, err := l.logsFile.WriteString("\n"); err != nil {
		return fmt.Errorf("failed to write newline: %w", err)
	}

	if err := l.writeHash(entry.SHA256); err != nil {
		return err
	}

	return nil
}

func (l *ImLogger) Info(message string) error {
	return l.Log(InfoLevel, message)
}

func (l *ImLogger) Error(message string) error {
	return l.Log(ErrorLevel, message)
}

func (l *ImLogger) Warning(message string) error {
	return l.Log(WarningLevel, message)
}

func (l *ImLogger) Event(message string) error {
	return l.Log(EventLevel, message)
}

func (l *ImLogger) Close() error {
	var errs []error

	if l.logsFile != nil {
		if err := l.logsFile.Close(); err != nil {
			errs = append(errs, err)
		}
	}

	if l.sumsFile != nil {
		if err := l.sumsFile.Close(); err != nil {
			errs = append(errs, err)
		}
	}

	if len(errs) > 0 {
		return fmt.Errorf("failed to close files: %v", errs)
	}

	return nil
}

func (l *ImLogger) GetLastHash() [32]byte {
	return l.lastSum
}