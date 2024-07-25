package alog

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strings"
	"time"
)

type Logger struct {
	writer io.Writer
	level  Level
	fields fields
}

type fields map[string]any

type Level int

const (
	Debug Level = iota
	Info
	Warning
	Error
	Panic
)

var levels = [5]string{"DBG", "INF", "WRN", "ERR", "CRT"}

type Option func(*Logger)

func WithLevel(level Level) Option {
	return func(l *Logger) {
		l.level = level
	}
}

func New(opts ...Option) Logger {
	l := Logger{}
	l.setDefaults()
	for _, opt := range opts {
		opt(&l)
	}
	return l
}

func (l *Logger) setDefaults() {
	l.writer = os.Stdout
	l.level = Info
}

func (l Logger) Debug(msg string, args ...any) error {
	if l.level > Debug {
		return nil
	}
	preparedMsg := l.prepareMsg(msg, Debug, args...)

	_, err := preparedMsg.WriteTo(l.writer)
	return err
}

func (l Logger) Info(msg string, args ...any) error {
	if l.level > Info {
		return nil
	}
	preparedMsg := l.prepareMsg(msg, Info, args...)

	_, err := preparedMsg.WriteTo(l.writer)
	return err
}

func (l Logger) Warning(msg string, args ...any) error {
	if l.level > Warning {
		return nil
	}
	preparedMsg := l.prepareMsg(msg, Warning, args...)

	_, err := preparedMsg.WriteTo(l.writer)
	return err
}

func (l Logger) Error(msg string, args ...any) error {
	if l.level > Error {
		return nil
	}
	preparedMsg := l.prepareMsg(msg, Error, args...)

	_, err := preparedMsg.WriteTo(l.writer)
	return err
}

func (l Logger) Panic(msg string, args ...any) error {
	preparedMsg := l.prepareMsg(msg, Panic, args...)

	_, err := preparedMsg.WriteTo(l.writer)
	return err
}

func (l Logger) prepareMsg(s string, lvl Level, args ...any) *strings.Reader {
	formattedMsg := fmt.Sprintf(s, args...)
	msg := fmt.Sprintf("%s | %s | %s", time.Now().Format(time.RFC3339), levels[lvl], formattedMsg)
	var fieldsMsg string

	if l.fields != nil {
		fields, err := json.Marshal(l.fields)
		if err == nil {
			fieldsMsg = string(fields)
		}
		msg = fmt.Sprintf("%s | %s", msg, fieldsMsg)
	}
	msg += "\n"

	return strings.NewReader(msg)
}

func (l Logger) WithField(key string, value any) Logger {
	if l.fields == nil {
		l.fields = make(fields)
	}
	newLogger := copyLogger(l)
	newLogger.fields[key] = value
	return newLogger
}
func copyLogger(l Logger) Logger {
	newFields := make(fields)
	for key, val := range l.fields {
		newFields[key] = val
	}
	logger := Logger{
		writer: l.writer,
		level:  l.level,
		fields: newFields,
	}
	return logger
}
