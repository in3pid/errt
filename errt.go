package errt

import (
	"bytes"
	"fmt"
	"runtime"
	"strings"
)

const MAXTRACE = 16

type ErrTrace struct {
	Err   error
	Trace []uintptr
}

func Trace(err error) error {
	trace := make([]uintptr, MAXTRACE)
	n := runtime.Callers(2, trace)
	return &ErrTrace{Err: err, Trace: trace[:n]}
}

func TraceDeferred(err error) error {
	trace := make([]uintptr, MAXTRACE)
	n := runtime.Callers(3, trace)
	return &ErrTrace{Err: err, Trace: trace[:n]}
}

func (e *ErrTrace) Error() string {
	g := func(s string) string {
		if n := strings.LastIndex(s, "/"); n >= 0 {
			s = s[n+1:]
		}
		return s
	}
	b := bytes.NewBuffer(nil)
	fmt.Fprintf(b, "%v\n", e.Err)
	for _, pc := range e.Trace {
		f := runtime.FuncForPC(pc)
		name := f.Name()
		file, line := f.FileLine(pc)
		fmt.Fprintf(b, "%s:%d %s\n", g(file), line, g(name))
	}
	return b.String()
}
