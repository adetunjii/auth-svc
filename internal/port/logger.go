package port

type AppLogger interface {
	Info(msg string, kv ...any)
	Error(msg string, err error)
	Fatal(msg string, err error)
}
