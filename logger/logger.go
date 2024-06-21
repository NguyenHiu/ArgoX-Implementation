package logger

import "fmt"

const (
	title_length = 35
	Reset        = "\033[0m"
	Red          = "\033[31m"
	Green        = "\033[32m"
	Yellow       = "\033[33m"
	Blue         = "\033[34m"
	Purple       = "\033[35m"
	Cyan         = "\033[36m"
	White        = "\033[37m"
)

type Logger struct {
	name string
}

func NewLogger(_name string) *Logger {
	return &Logger{
		name: _name,
	}
}

func (l *Logger) ChangeName(_newName string) {
	l.name = _newName
}

func (l *Logger) Debug(format string, a ...any) {
	fmt.Printf("%-*s: %v", title_length, fmt.Sprintf("[%sDEBUG%s] - %s", Blue, Reset, l.name), fmt.Sprintf(format, a...))
}

func (l *Logger) Info(format string, a ...any) {
	fmt.Printf("%-*s: %v", title_length, fmt.Sprintf("[%sINFO%s] - %s", Green, Reset, l.name), fmt.Sprintf(format, a...))
}

func (l *Logger) Error(format string, a ...any) {
	fmt.Printf("%-*s: %v", title_length, fmt.Sprintf("[%sERROR%s] - %s", Red, Reset, l.name), fmt.Sprintf(format, a...))
}
