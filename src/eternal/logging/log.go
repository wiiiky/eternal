package logging

import (
	log "github.com/sirupsen/logrus"
	"io"
	"os"
)

func Start(format, level, output string) {
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

	var f io.Writer
	var err error
	if output == "stdout" {
		f = os.Stdout
	} else if output == "stderr" {
		f = os.Stderr
	} else {
		f, err = os.OpenFile(output, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	}
	if err != nil {
		panic(err)
	}
	log.SetOutput(f)

	log.AddHook(ContextHook{})
}
