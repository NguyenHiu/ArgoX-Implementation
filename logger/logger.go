package logger

import "fmt"

const (
	title_length = 35
	Reset        = "\033[0m"
	Bold         = "\033[1m"
	Italic       = "\033[3m"
	Red          = "\033[31m"
	Green        = "\033[32m"
	Yellow       = "\033[33m"
	Blue         = "\033[34m"
	Purple       = "\033[35m"
	Cyan         = "\033[36m"
	White        = "\033[37m"
	Magenta      = "\033[35m"
	LightGray    = "\033[37m"
	DarkGray     = "\033[90m"
	LightRed     = "\033[91m"
	LightGreen   = "\033[92m"
	LightYellow  = "\033[93m"
	LightBlue    = "\033[94m"
	LightMagenta = "\033[95m"
	LightCyan    = "\033[96m"
	LightWhite   = "\033[97m"
	None         = ""
)

type Logger struct {
	name         string
	color        string
	format       string
	formatedName string
}

func NewLogger(_name string, _color string, _format string) *Logger {
	_reset := Reset
	if _color == "" && _format == "" {
		_reset = ""
	}
	return &Logger{
		name:         _name,
		color:        _color,
		format:       _format,
		formatedName: fmt.Sprintf("%s%s%v%s", _format, _color, _name, _reset),
	}
}

func (l *Logger) ChangeName(_newName string, _color string, _format string) {
	l.name = _newName
	l.color = _color
	l.format = _format
	_reset := Reset
	if _color == "" && _format == "" {
		_reset = ""
	}
	l.formatedName = fmt.Sprintf("%s%s%v%s", l.format, l.color, l.name, _reset)
}

func (l *Logger) Debug(format string, a ...any) {
	fmt.Printf("%-*s: %v", title_length+len(l.formatedName)-len(l.name), fmt.Sprintf("[%sDEBUG%s] - %s", Blue, Reset, l.formatedName), fmt.Sprintf(format, a...))
}

func (l *Logger) Info(format string, a ...any) {
	fmt.Printf("%-*s: %v", title_length+len(l.formatedName)-len(l.name), fmt.Sprintf("[%sINFO%s] - %s", Green, Reset, l.formatedName), fmt.Sprintf(format, a...))
}

func (l *Logger) Error(format string, a ...any) {
	fmt.Printf("%-*s: %v", title_length+len(l.formatedName)-len(l.name), fmt.Sprintf("[%sERROR%s] - %s", Red, Reset, l.formatedName), fmt.Sprintf(format, a...))
}
