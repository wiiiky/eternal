package logging

import (
	log "github.com/sirupsen/logrus"
	"path"
	"runtime"
)

type ContextHook struct{}

func (hook ContextHook) Levels() []log.Level {
	return log.AllLevels
}

func (hook ContextHook) Fire(entry *log.Entry) error {
	if pc, file, line, ok := runtime.Caller(9); ok {
		funcName := runtime.FuncForPC(pc).Name()

		entry.Data["file"] = path.Base(file)
		entry.Data["func"] = path.Base(funcName)
		entry.Data["line"] = line
	}

	return nil
}
