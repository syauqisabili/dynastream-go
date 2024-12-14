package pkg

import (
	_ "embed"
	"fmt"
	"io"
	"os"
	"runtime"
	"sync"

	"github.com/charmbracelet/log"
	"github.com/joho/godotenv"
	"gopkg.in/natefinch/lumberjack.v2"
)

// Logger is a wrapper around the charmbracelet/log.Logger
var (
	logger *log.Logger
	once   sync.Once
)

// logFile is defined at the package level to manage its lifecycle
var LogFile *lumberjack.Logger

// init initializes the logger instance.
func init() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Failed to read .env")
		os.Exit(1)
	}

	// Select runtime
	var applicationDir string
	runtime := runtime.GOOS
	switch runtime {
	case "windows":
		applicationDir = os.Getenv("DATA_DIR_WIN")
	case "linux":
		dir, err := os.UserHomeDir()
		if err != nil {
			log.Fatal("Failed to find home directory")
			os.Exit(1)
		}

		applicationDir = dir + "/data"
	}
	// Create logs directory if it doesn't exist
	dir := fmt.Sprintf("%s/%s", applicationDir, os.Getenv("APPLICATION_NAME"))
	if err := os.MkdirAll(dir, os.ModePerm); err != nil {
		log.Fatal("Error creating log directory:", err)
		os.Exit(1)
	}

	pathFile := fmt.Sprintf("%s/%s/%s", applicationDir, os.Getenv("APPLICATION_NAME"), os.Getenv("FILENAME_LOG"))
	LogFile = &lumberjack.Logger{
		Filename:   pathFile, // Log file name
		MaxSize:    10,       // Maximum size in megabytes before rotation
		MaxBackups: 30,       // Maximum number of backup log files to keep
		MaxAge:     30,       // Maximum number of days to retain old log files
		Compress:   true,     // Compress old log files
	}

	once.Do(func() {
		multiWriter := io.MultiWriter(LogFile, os.Stdout)
		logger = log.New(multiWriter)
		logger.SetFormatter(log.TextFormatter)
		logger.SetReportCaller(false)
		logger.SetReportTimestamp(true)

	})
}

func LogInfo(msg interface{}, keyvals ...interface{}) {
	logger.Info(msg, keyvals...)
}

func LogDebug(msg interface{}, keyvals ...interface{}) {
	logger.Debug(msg, keyvals...)
}

func LogWarn(msg interface{}, keyvals ...interface{}) {
	logger.Warn(msg, keyvals...)
}

func LogError(msg interface{}, keyvals ...interface{}) {
	logger.Error(msg, keyvals...)
}

func LogFatal(msg interface{}, keyvals ...interface{}) {
	logger.Fatal(msg, keyvals...)
}
