package logging

import (
	"fmt"
	"os"
)

type StdOutLog struct {
	level  LogLevel
	stdout *os.File
}

func (s *StdOutLog) SetLevel(l LogLevel) {
	s.level = l
}

func (s *StdOutLog) Logf(level LogLevel, m string, format ...interface{}) {
	if s.level < level {
		return
	}
	if len(format) > 0 {
		fmt.Printf(m, format)
		return
	}
	fmt.Println(m)
}
