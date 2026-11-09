// Package log provides the application's shared logger.
package log

import (
	"fmt"
	"os"
	"strings"

	"github.com/sirupsen/logrus"
)

type Fields = logrus.Fields

var Logger = logrus.New()

func init() {
	if err := Setup(); err != nil {
		panic(err)
	}
}

func Setup() error {
	levelName := strings.TrimSpace(os.Getenv("LOG_LEVEL"))
	if levelName == "" {
		levelName = logrus.InfoLevel.String()
	}

	level, err := logrus.ParseLevel(strings.ToLower(levelName))
	if err != nil {
		return fmt.Errorf("parse LOG_LEVEL %q: %w", levelName, err)
	}

	Logger.SetOutput(os.Stderr)
	Logger.SetFormatter(&logrus.JSONFormatter{})
	Logger.SetLevel(level)

	return nil
}

func WithFields(fields Fields) *logrus.Entry {
	return Logger.WithFields(fields)
}

func Trace(args ...any) {
	Logger.Trace(args...)
}
func Debug(args ...any) {
	Logger.Debug(args...)
}

func Info(args ...any) {
	Logger.Info(args...)
}

func Warn(args ...any) {
	Logger.Warn(args...)
}

func Error(args ...any) {
	Logger.Error(args...)
}

func Fatal(args ...any) {
	Logger.Fatal(args...)
}

func Panic(args ...any) {
	Logger.Panic(args...)
}

func Tracef(format string, args ...any) {
	Logger.Tracef(format, args...)
}

func Debugf(format string, args ...any) {
	Logger.Debugf(format, args...)
}

func Infof(format string, args ...any) {
	Logger.Infof(format, args...)
}

func Warnf(format string, args ...any) {
	Logger.Warnf(format, args...)
}

func Errorf(format string, args ...any) {
	Logger.Errorf(format, args...)
}

func Fatalf(format string, args ...any) {
	Logger.Fatalf(format, args...)
}

func Panicf(format string, args ...any) {
	Logger.Panicf(format, args...)
}
