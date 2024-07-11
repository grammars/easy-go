package logger

import (
	"github.com/grammars/easy-go/file"
	"gopkg.in/natefinch/lumberjack.v2"
	"io"
	"log/slog"
	"os"
	"path/filepath"
)

type Option struct {
	ConsoleEnabled bool
	FileEnabled    bool
	JsonMode       bool
	LogLevel       slog.Level

	// Filename is the file to write logs to.  Backup log files will be retained
	// in the same directory.  It uses <processname>-lumberjack.log in
	// os.TempDir() if empty.
	Filename string

	// MaxSize is the maximum size in megabytes of the log file before it gets
	// rotated. It defaults to 100 megabytes.
	MaxSize int

	// MaxAge is the maximum number of days to retain old log files based on the
	// timestamp encoded in their filename.  Note that a day is defined as 24
	// hours and may not exactly correspond to calendar days due to daylight
	// savings, leap seconds, etc. The default is not to remove old log files
	// based on age.
	MaxAge int

	// MaxBackups is the maximum number of old log files to retain.  The default
	// is to retain all old log files (though MaxAge may still cause them to get
	// deleted.)
	MaxBackups int

	// LocalTime determines if the time used for formatting the timestamps in
	// backup files is the computer's local time.  The default is to use UTC
	// time.
	LocalTime bool

	// Compress determines if the rotated log files should be compressed
	// using gzip. The default is not to perform compression.
	Compress bool
}

type ComboWriter struct {
	consoleWriter io.Writer
	fileWriter    io.Writer
}

func (w *ComboWriter) Write(p []byte) (n int, err error) {
	cw, err := w.consoleWriter.Write(p)
	if err != nil {
		return cw, err
	}
	return w.fileWriter.Write(p)
}

func CreateOption() *Option {
	return &Option{ConsoleEnabled: true, FileEnabled: false, JsonMode: true, LogLevel: slog.LevelInfo,
		MaxSize: 2, MaxBackups: 3, MaxAge: 7, LocalTime: true, Compress: false}
}

func (option *Option) Setup() {
	logLevel := slog.LevelDebug
	if option.Filename == "" {
		option.Filename = filepath.Join(file.GetExeDir(), ".logs", "app.log")
	} else {
		option.FileEnabled = true
	}

	var fileWriter io.Writer
	if option.FileEnabled {
		fileWriter = &lumberjack.Logger{
			Filename:   option.Filename,
			MaxSize:    option.MaxSize, // megabytes
			MaxBackups: option.MaxBackups,
			MaxAge:     option.MaxAge, // days
			Compress:   option.Compress,
		}
	}
	var consoleWriter io.Writer
	if option.ConsoleEnabled {
		consoleWriter = os.Stdout
	}

	var writer io.Writer
	if option.FileEnabled && option.ConsoleEnabled {
		writer = &ComboWriter{fileWriter: fileWriter, consoleWriter: consoleWriter}
	} else if option.FileEnabled {
		writer = fileWriter
	} else {
		writer = consoleWriter
	}

	opts := &slog.HandlerOptions{Level: logLevel}
	var logger *slog.Logger
	if option.JsonMode {
		logger = slog.New(slog.NewJSONHandler(writer, opts))
	} else {
		logger = slog.New(slog.NewTextHandler(writer, opts))
	}
	slog.SetDefault(logger)
	slog.Info("Logger设置完毕", "选项", option)
}
