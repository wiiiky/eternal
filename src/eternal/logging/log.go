package logging

import (
	log "github.com/sirupsen/logrus"
	"io"
	"os"
)

func OpenLogFile(filename string) (io.Writer, error) {
	var f io.Writer
	var err error
	if filename == "stdout" {
		f = os.Stdout
	} else if filename == "stderr" {
		f = os.Stderr
	} else {
		f, err = os.OpenFile(filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	}
	return f, err
}

func Init(format, level, output string) {
	allFormatters := map[string]log.Formatter{
		"json": &log.JSONFormatter{},
		"text": &log.TextFormatter{},
	}
	formatter, ok := allFormatters[format]
	if !ok {
		panic("Invalid log format. (json/text)")
	}
	log.SetFormatter(formatter)

	allLevels := map[string]log.Level{
		"panic": log.PanicLevel,
		"fatal": log.FatalLevel,
		"error": log.ErrorLevel,
		"warn":  log.WarnLevel,
		"info":  log.InfoLevel,
		"debug": log.DebugLevel,
	}
	lvel, ok := allLevels[level]
	if !ok {
		panic("Invalid log level. (debug/info/warn/error/fatal/panic)")
	}
	log.SetLevel(lvel)

	f, err := OpenLogFile(output)
	if err != nil {
		panic(err)
	}
	log.SetOutput(f)

	log.AddHook(ContextHook{})
}
