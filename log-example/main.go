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

// Fire write log ke file
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

	// create log file
	logFile := "error.log"
	file, err := os.OpenFile(logFile, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		logger.Fatalf("Failed to create file hook: %v", err)
	}
	defer file.Close()

	// Added hook for error log
	logger.AddHook(&FileHook{
		writer:    file,
		logLevels: []logrus.Level{logrus.ErrorLevel, logrus.FatalLevel, logrus.PanicLevel},
	})

	logger.SetReportCaller(true)
	logger.SetFormatter(&logrus.TextFormatter{
		FullTimestamp: true,
		DisableColors: true,
	})

	// Sets the output for the terminal
	logger.SetOutput(os.Stdout)

	// examples
	logger.Info("This is an info message")
	logger.Warn("This is a warning message")
	logger.Error("This is an error message")
}
