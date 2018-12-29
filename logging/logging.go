package logging

type Type int
type LogLevel int

const ( // Log facilities
	Stdout Type = iota
	File
)

const (
	Stage LogLevel = iota
	Error
	Database
	HTTP
	All
)

type Log interface {
	SetLevel(LogLevel)
	Logf(LogLevel, string, ...interface{})
}

var Journal Log

func NewLogger(t Type, l LogLevel) error {

	switch t {
	case Stdout:
		Journal = new(StdOutLog)
		Journal.SetLevel(l)
	}

	return nil
}
