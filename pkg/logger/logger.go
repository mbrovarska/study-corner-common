package log

// Field is a simple key/value wrapper that does not depend on zap
type Field struct {
    Key   string
    Value any
}

// Logger is the interface that microservices can depend on from core code
type Logger interface {
    Debug(msg string, fields ...Field)
    Info(msg string, fields ...Field)
    Warn(msg string, fields ...Field)
    Error(msg string, fields ...Field)
    With(fields ...Field) Logger
}

// Helper constructors for fields so core code looks nice
func String(key, value string) Field {
    return Field{Key: key, Value: value}
}

func Int64(key string, value int64) Field {
    return Field{Key: key, Value: value}
}


