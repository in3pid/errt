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

func trace(p []uintptr, adjust int, filter Filter) (n int) {
	for ; n < len(p); n++ {
		pc, file, line, ok := runtime.Caller(n + adjust)
		if !ok {
			break
		}
		name := runtime.FuncForPC(pc).Name()
		if !filter(name, file, line) {
			break
		}
		p[n] = pc
	}
	return n
}

type Filter func(name, file string, line int) bool

func PackageFilter() Filter {
	var prefix *string
	return func(name, _ string, _ int) bool {
		n := 0
		if n = strings.LastIndex(name, "/"); 0 <= n {
			name = name[n+1:]
		}
		if prefix == nil {
			if n := strings.Index(name, "."); 0 <= n {
				name = name[:n+1]
			}
			prefix = &name
			return true
		}
		return strings.HasPrefix(name, *prefix)
	}
}

func Trace(err error) error {
	p := make([]uintptr, MAXTRACE)
	n := trace(p, 2, PackageFilter())
	return &ErrTrace{Err: err, Trace: p[:n]}
}

func TraceDeferred(err error) error {
	p := make([]uintptr, MAXTRACE)
	n := trace(p, 3, PackageFilter())
	return &ErrTrace{Err: err, Trace: p[:n]}
}

func TraceAll(err error) error {
	trace := make([]uintptr, MAXTRACE)
	n := runtime.Callers(2, trace)
	return &ErrTrace{Err: err, Trace: trace[:n]}
}

func TraceAllDeferred(err error) error {
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
