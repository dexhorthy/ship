package logger

import (
	"fmt"
	"io"
	"strings"
	"time"

	"github.com/go-kit/kit/log"
	"github.com/go-logfmt/logfmt"
)

var _ log.Logger = &prettyLogger{}

type prettyLogger struct {
	w io.Writer
}

var DefaultTimestampPretty = log.TimestampFormat(
	func() time.Time { return time.Now().UTC() },
	time.RFC3339,
)

func (p *prettyLogger) Log(keyvals ...interface{}) error {
	event := ""
	msg := ""
	var sstruct string
	var method string
	var ts string
	var level string
	var next *string
	var unknown []interface{}

	for index, keyval := range keyvals {
		if index%2 == 1 {
			if v, ok := keyval.(log.Valuer); ok {
				keyval = v()
			}
			if next != nil {
				*next = fmt.Sprintf("%s", keyval)
			} else {
				unknown = append(unknown, fmt.Sprintf("%s", keyval))
			}
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
		case "level":
			next = &level
			continue
		}

		next = nil
		unknown = append(unknown, keyval)
	}

	eventColor := Blue
	switch level {
	case "error":
		eventColor = Red
		break
	case "warn":
		eventColor = Yellow
		break
	case "info":
		eventColor = Green
		break
	case "debug":
		fallthrough
	default:
		eventColor = Blue
	}

	var caller string
	if method != "" && sstruct != "" {
		caller = fmt.Sprintf("%s [%s.%s]%s", fgColorBytes[White], method, sstruct, resetColorBytes)
	} else if method != "" {
		caller = fmt.Sprintf("%s [%s]%s", fgColorBytes[White], method, resetColorBytes)
	}

	if event != "" {
		event = fmt.Sprintf(" %s%25s%s", fgColorBytes[eventColor], event, resetColorBytes)
	}

	if ts != "" {
		ts = fmt.Sprintf("%s%s%s", resetColorBytes, ts, resetColorBytes)
	}

	unknownstr := ""
	if len(unknown) != 0 {
		bytes, e := logfmt.MarshalKeyvals(unknown)
		if e != nil {
			bytes = []byte(fmt.Sprintf("(+%d unknown)", len(unknown)))
		}
		unknownstr = fmt.Sprintf(" %s(%s)%s", fgColorBytes[Gray], bytes, resetColorBytes)
	}

	levelstr := fmt.Sprintf("%s%s%s", fgColorBytes[eventColor], strings.ToUpper(level), resetColorBytes)

	_, err := p.w.Write([]byte(fmt.Sprintf("%s %5s%25s%s %s%s\n", ts, levelstr, event, caller, msg, unknownstr)))
	return err
}
