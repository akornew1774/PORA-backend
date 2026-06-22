// logger - пакет, обеспечивающий кастомное логирование
package logger

import (
	"io"
	"os"

	log "github.com/sirupsen/logrus"
	"gopkg.in/natefinch/lumberjack.v2"
)

// Log - кастомный логгер для приложения
var Log *log.Logger

// InitLogger инициализирует новый кастомный логгер,
// позволяющитй логировать и в консоль, и в отдельный файл
func InitLogger() {
	Log = log.New()

	_, err := os.Stat("logs")
	if os.IsNotExist(err) {
		os.Mkdir("logs", 0755)
	}

	logFile := &lumberjack.Logger{
		Filename:   "logs/pora.log",
		MaxSize:    20,
		MaxBackups: 5,
		MaxAge:     30,
		Compress:   true,
	}

	Log.SetOutput(io.MultiWriter(os.Stdout, logFile))

	Log.SetFormatter(&log.TextFormatter{FullTimestamp: true})
	Log.SetLevel(log.InfoLevel)
}

// InitLoggerTest создает кастомный логгер для тестов,
// который отключает все логи
func InitLoggerTest() {
	Log = log.New()
	Log.SetOutput(io.Discard)
	Log.SetLevel(log.InfoLevel)
}
