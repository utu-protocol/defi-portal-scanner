package utils

import (
	"fmt"

	"github.com/sirupsen/logrus"
)

// EmojiLogFormatter struct for emoji log formatter
type EmojiLogFormatter struct {
}

// Format format a log message
func (f *EmojiLogFormatter) Format(entry *logrus.Entry) ([]byte, error) {
	format := "%s\n"
	switch entry.Level {
	case logrus.DebugLevel:
		format = "🔍 %s\n"
	case logrus.ErrorLevel:
		format = "❌  %s\n"
	case logrus.FatalLevel:
		format = "☠️  %s\n"
	case logrus.WarnLevel:
		format = "⚠️  %s\n"
	}

	if isSuccess, found := entry.Data["success"]; found && isSuccess.(bool) {
		format = "✅ %s\n"
	}

	line := fmt.Sprintf(format, entry.Message)
	return []byte(line), nil
}
