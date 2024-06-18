package main

import (
	"os"

	"github.com/sirupsen/logrus"
)

// FileHook struct for hook file
type FileHook struct {
	writer    *os.File
	logLevels []logrus.Level
}

// NewFileHook create new FileHook
func NewFileHook(logFile string, levels []logrus.Level) (*FileHook, error) {
	file, err := os.OpenFile(logFile, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		return nil, err
	}
	return &FileHook{
		writer:    file,
		logLevels: levels,
	}, nil
}

// Fire menulis log ke file
func (hook *FileHook) Fire(entry *logrus.Entry) error {
	line, err := entry.String()
	if err != nil {
		return err
	}
	_, err = hook.writer.Write([]byte(line))
	return err
}

// Levels returns the levels logged by this hook
func (hook *FileHook) Levels() []logrus.Level {
	return hook.logLevels
}

func main() {
	// init logrus
	logger := logrus.New()

	// Added hook for error log
	fileHook, err := NewFileHook("error.log", []logrus.Level{logrus.ErrorLevel, logrus.FatalLevel, logrus.PanicLevel})
	if err != nil {
		logger.Fatalf("Failed to create file hook: %v", err)
	}
	logger.AddHook(fileHook)

	// Sets the output for the terminal
	logger.SetOutput(os.Stdout)

	// examples
	logger.Info("This is an info message")
	logger.Warn("This is a warning message")
	logger.Error("This is an error message")
}
