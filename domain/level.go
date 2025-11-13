package domain

type Level int8

const (
	DebugLevel Level = iota - 1
	InfoLevel
	WarnLevel
	ErrorLevel
)

type LevelAware interface {
	LevelEnabled(level Level) bool
}
