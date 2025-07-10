package u

import (
	"fmt"
	"runtime"
	"strings"
)

type Err struct {
	Message string
	Stack   []string
}

func (e *Err) Error() string {
	return e.Message
}

func Errorf(format string, args ...interface{}) *Err {
	return Error(fmt.Sprintf(format, args...))
}

func Error(err string) *Err {
	callStacks := make([]string, 0)
	for i := 1; i < 50; i++ {
		_, file, line, ok := runtime.Caller(i)
		if !ok {
			break
		}
		if strings.Contains(file, "/go/src/") {
			continue
		}
		if strings.Contains(file, "/ssgo/log") {
			continue
		}
		if strings.Contains(file, "/ssgo") {
			file = "**" + file[strings.LastIndex(file, "/ssgo"):]
		}
		callStacks = append(callStacks, fmt.Sprintf("%s:%d", file, line))
	}
	return &Err{
		Message: err,
		Stack:   callStacks,
	}
}
