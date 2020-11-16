package logging

import (
	"context"
	"fmt"
	"io"
	logging "log"
	"os"
	"runtime"
	"strings"

	"github.com/sirupsen/logrus"
)

var (
	log *logrus.Logger
    logFile *os.File = nil
)

func Close() {
	if logFile != nil {
		err := logFile.Close()
		if err != nil {
			Error(nil,err,"%v", "Unable to properly close the logfile")
		}
	}
}

func GetLogger() *logrus.Logger {
	return log
}

func Init(path string, logLevel string, withColors bool) {
	f, err := os.OpenFile(path, os.O_WRONLY | os.O_CREATE | os.O_APPEND, 0755)

	if err != nil {
		logging.Fatalf("error opening file: %v", err)
	}

	log = logrus.New()
	log.SetFormatter(&logrus.TextFormatter{
		ForceColors: withColors,
		FullTimestamp: true,
		PadLevelText: true,
	})
	mw := io.MultiWriter(os.Stdout, f)
	level, err := logrus.ParseLevel(logLevel)

	if err != nil {
		level = logrus.InfoLevel
		Error(nil,err,"%v", "Unable to parse log level (none of fatal, error, warning, info, debug or trace), defaulting to info level")
	}

	log.SetLevel(level)
	log.SetOutput(mw)
}

func init() {
	log = logrus.New()
	log.SetFormatter(&logrus.TextFormatter{
		FullTimestamp: true,
		PadLevelText: true,
	})
	log.SetLevel(logrus.InfoLevel)
}

func Trace(ctx context.Context, format string, msg ...interface{}) {
	getLogEntry(ctx).Tracef(format, msg...)
}

func getLogEntry(ctx context.Context) *logrus.Entry {
	if ctx == nil {
		ctx = context.Background()
	}

	if pc, file, line, ok := runtime.Caller(2); ok {
		file = file[strings.LastIndex(file, "/")+1:]
		funcName := runtime.FuncForPC(pc).Name()
		return log.WithFields(
			logrus.Fields{
				"src": fmt.Sprintf("%s:%d (%s)", file, line, funcName),
			}).WithContext(ctx)
	}

	return logrus.NewEntry(log)
}

func Debug(ctx context.Context, format string, msg ...interface{}) {
	getLogEntry(ctx).Debugf(format, msg...)
}

func Info(ctx context.Context, format string, msg ...interface{}) {
	getLogEntry(ctx).Infof(format, msg...)
}

func Warn(ctx context.Context, format string, msg ...interface{}) {
	getLogEntry(ctx).Warnf(format, msg...)
}

func Error(ctx context.Context, err error, format string, msg ...interface{}) {
	getLogEntry(ctx).WithError(err).Errorf(format, msg...)
}
