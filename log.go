package utils

import (
	"bytes"
	"fmt"
	"os"
	"runtime"
	"text/template"
	"time"

	"github.com/sirupsen/logrus"
)

type templateFormatter struct {
	template        *template.Template
	timestampFormat string
}

func (f *templateFormatter) Format(entry *logrus.Entry) ([]byte, error) {

	r := entry.Message
	m := make(map[string]interface{})

	for k, v := range entry.Data {
		switch v := v.(type) {
		case error:
			m[k] = v.Error()
		default:
			m[k] = v
		}
	}

	m["msg"] = entry.Message
	m["time"] = entry.Time.Format(f.timestampFormat)
	m["level"] = entry.Level.String()

	var err error

	if f.template != nil {

		var b bytes.Buffer
		err = f.template.Execute(&b, m)
		if err == nil {

			r = fmt.Sprintf("%s\n", b.String())
		}
	}

	return []byte(r), err
}

type Log struct {
	CallInfo bool
}

var log = Log{}

func GetLog() *Log {

	return &log
}

func (log *Log) trace(offset int) logrus.Fields {

	if log.CallInfo {

		pc := make([]uintptr, 15)
		n := runtime.Callers(offset, pc)
		frames := runtime.CallersFrames(pc[:n])
		frame, _ := frames.Next()
		return logrus.Fields{
			"file": frame.File,
			"line": frame.Line,
			"func": frame.Function,
		}
	} else {

		return logrus.Fields{}
	}
}

func prepare(message string, args ...interface{}) string {

	if len(args) > 0 {

		return fmt.Sprintf(message, args...)
	} else {

		return message
	}
}

func exists(level logrus.Level, obj interface{}, args ...interface{}) (bool, string) {

	if obj == nil {

		return false, ""
	}

	message := ""

	switch v := obj.(type) {
	case error:
		message = v.Error()
	case string:
		message = v
	default:
		message = "not implemented"
	}

	flag := message != "" && logrus.IsLevelEnabled(level)
	if flag {
		message = prepare(message, args...)
	}
	return flag, message
}

func (log *Log) Info(obj interface{}, args ...interface{}) int64 {

	if exists, message := exists(logrus.InfoLevel, obj, args...); exists {

		logrus.WithFields(log.trace(3)).Infoln(message)
	}
	return time.Now().UnixNano()
}

func (log *Log) Warn(obj interface{}, args ...interface{}) int64 {

	if exists, message := exists(logrus.WarnLevel, obj, args...); exists {

		logrus.WithFields(log.trace(3)).Warnln(message)
	}
	return time.Now().UnixNano()
}

func (log *Log) Error(obj interface{}, args ...interface{}) int64 {

	if exists, message := exists(logrus.ErrorLevel, obj, args...); exists {

		logrus.WithFields(log.trace(3)).Errorln(message)
	}
	return time.Now().UnixNano()
}

func (log *Log) Debug(obj interface{}, args ...interface{}) int64 {

	if exists, message := exists(logrus.DebugLevel, obj, args...); exists {

		logrus.WithFields(log.trace(3)).Debugln(message)
	}
	return time.Now().UnixNano()
}

func (log *Log) Panic(obj interface{}, args ...interface{}) {

	if exists, message := exists(logrus.PanicLevel, obj, args...); exists {

		logrus.WithFields(log.trace(3)).Panicln(message)
	}
}

func (log *Log) Init(format string, level string, templ string) {

	switch format {
	case "json":
		logrus.SetFormatter(&logrus.JSONFormatter{})
	case "text":
		logrus.SetFormatter(&logrus.TextFormatter{})
	case "stdout":
		t, err := template.New("").Parse(templ)
		if err != nil {
			logrus.Panic(err)
		}
		logrus.SetFormatter(&templateFormatter{template: t, timestampFormat: time.RFC3339})
	}

	switch level {
	case "info":
		logrus.SetLevel(logrus.InfoLevel)
	case "error":
		logrus.SetLevel(logrus.ErrorLevel)
	case "panic":
		logrus.SetLevel(logrus.PanicLevel)
	case "warn":
		logrus.SetLevel(logrus.WarnLevel)
	case "debug":
		logrus.SetLevel(logrus.DebugLevel)

	}

	logrus.SetOutput(os.Stdout)
}
