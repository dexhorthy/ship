package logger

import (
	"fmt"
	golog "log"
	"os"
	"path"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"github.com/go-stack/stack"
	multierror "github.com/hashicorp/go-multierror"
	"github.com/mitchellh/cli"
	"github.com/replicatedhq/ship/pkg/constants"
	"github.com/replicatedhq/ship/pkg/ui"
	"github.com/spf13/afero"
	"github.com/spf13/viper"
)

type compositeLogger struct {
	loggers []log.Logger
}

var _ log.Logger = &prettyLogger{}

type prettyLogger struct {
	UI cli.Ui
}

func (p *prettyLogger) Log(keyvals ...interface{}) error {
	event := ""
	msg := ""
	var sstruct string
	var method string
	var ts string
	var next *string
	var devnull string
	var unknown int

	for index, keyval := range keyvals {
		if index%2 == 1 {
			if v, ok := keyval.(log.Valuer); ok {
				keyval = v()
			}
			*next = fmt.Sprintf("%s", keyval)
			continue
		}

		switch keyval {
		case "event":
			next = &event
			continue
		case "struct":
			next = &sstruct
			continue
		case "method":
			next = &method
			continue
		case "message":
			fallthrough
		case "msg":
			next = &method
			continue
		case "ts":
			next = &ts
			continue
		}

		next = &devnull // drop unknown keys for now
		unknown += 1
	}

	p.UI.Output(fmt.Sprintf("%s [%s.%s]: %s %s (+%d unknown)", ts, sstruct, method, event, msg, unknown))
	return nil
}

func (c *compositeLogger) Log(keyvals ...interface{}) error {
	var multiErr *multierror.Error
	for _, logger := range c.loggers {
		multiErr = multierror.Append(multiErr, logger.Log(keyvals...))
	}
	return multiErr.ErrorOrNil()
}

// New builds a logger from env using viper
func New(v *viper.Viper, fs afero.Afero) log.Logger {

	fullPathCaller := pathCaller(6)
	var stdoutLogger log.Logger
	stdoutLogger = withFormat(viper.GetString("log-format"), v)
	stdoutLogger = log.With(stdoutLogger, "ts", log.DefaultTimestampUTC)
	stdoutLogger = log.With(stdoutLogger, "caller", fullPathCaller)
	stdoutLogger = withLevel(stdoutLogger, v.GetString("log-level"))

	debugLogFile := path.Join(constants.ShipPathInternalLog)
	var debugLogger log.Logger
	err := fs.RemoveAll(debugLogFile)
	if err != nil {
		level.Warn(stdoutLogger).Log("msg", "failed to remove existing debug log file", "path", debugLogFile, "error", err)
		golog.SetOutput(log.NewStdlibAdapter(level.Debug(stdoutLogger)))
		return stdoutLogger
	}
	debugLogWriter, err := fs.Create(debugLogFile)
	if err != nil {
		level.Warn(stdoutLogger).Log("msg", "failed to initialize debug log file", "path", debugLogFile, "error", err)
		golog.SetOutput(log.NewStdlibAdapter(level.Debug(stdoutLogger)))
		return stdoutLogger
	}

	debugLogger = log.NewJSONLogger(debugLogWriter)
	debugLogger = log.With(debugLogger, "ts", log.DefaultTimestampUTC)
	debugLogger = log.With(debugLogger, "caller", fullPathCaller)
	debugLogger = withLevel(debugLogger, "debug")

	realLogger := &compositeLogger{
		loggers: []log.Logger{
			stdoutLogger,
			debugLogger,
		},
	}

	golog.SetOutput(log.NewStdlibAdapter(level.Debug(realLogger)))
	return realLogger
}

func withFormat(format string, v *viper.Viper) log.Logger {
	switch format {
	case "json":
		return log.NewJSONLogger(os.Stdout)
	case "logfmt":
		return log.NewLogfmtLogger(os.Stdout)
	case "pretty":
		return &prettyLogger{ui.FromViper(v)}
	default:
		return log.NewLogfmtLogger(os.Stdout)
	}

}

func withLevel(logger log.Logger, lvl string) log.Logger {
	switch lvl {
	case "debug":
		return level.NewFilter(logger, level.AllowDebug())
	case "info":
		return level.NewFilter(logger, level.AllowInfo())
	case "warn":
		return level.NewFilter(logger, level.AllowWarn())
	case "error":
		return level.NewFilter(logger, level.AllowError())
	case "off":
		return level.NewFilter(logger, level.AllowNone())
	default:
		logger.Log("msg", "Unknown log level, using debug", "received", lvl)
		return level.NewFilter(logger, level.AllowDebug())
	}
}

func pathCaller(depth int) log.Valuer {
	return func() interface{} {
		return fmt.Sprintf("%+s", stack.Caller(depth))
	}
}
