package logx

import (
	"fmt"
	"runtime"
	"strings"

	"github.com/sirupsen/logrus"
)

type Caller struct{}

func (c Caller) Levels() []logrus.Level {
	return logrus.AllLevels
}

func (c Caller) Fire(entry *logrus.Entry) error {
	frame := getCaller()
	if frame == nil {
		return nil
	}
	entry.Data["file"] = fmt.Sprintf(" %s:%d ", frame.File, frame.Line)
	entry.Data["func"] = frame.Function
	entry.Data["trace_id"] = ""
	if entry.Context != nil {
		entry.Data["trace_id"] = fmt.Sprintf("%v", entry.Context.Value("trace_id"))
	}
	return nil
}

// getCaller retrieves the name of the first non-logrus calling function
func getCaller() *runtime.Frame {
	pcs := make([]uintptr, 25)
	depth := runtime.Callers(4, pcs)
	frames := runtime.CallersFrames(pcs[:depth])
	logrusPackage := "github.com/sirupsen/logrus"
	restyxPackage := "resty"
	for f, again := frames.Next(); again; f, again = frames.Next() {
		pkg := getPackageName(f.Function)
		if !strings.Contains(pkg, logrusPackage) && !strings.Contains(pkg, restyxPackage) {
			return &f //nolint:scopelint
		}
	}
	return nil
}

// getPackageName reduces a fully qualified function name to the package name
// There really ought to be to be a better way...
func getPackageName(f string) string {
	for {
		lastPeriod := strings.LastIndex(f, ".")
		lastSlash := strings.LastIndex(f, "/")
		if lastPeriod > lastSlash {
			f = f[:lastPeriod]
		} else {
			break
		}
	}

	return f
}
