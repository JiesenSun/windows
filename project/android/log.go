package client

var (
	g_logger Logger = nil
)

type Logger interface {
	Debug(logMsg string)
	Info(logMsg string)
	Warn(logMsg string)
	Fatal(logMsg string)
}

type consoleLog struct{}

func (*consoleLog) Debug(logMsg string) {
	println(logMsg)
}

func (*consoleLog) Info(logMsg string) {
	println(logMsg)
}

func (*consoleLog) Warn(logMsg string) {
	println(logMsg)
}

func (*consoleLog) Fatal(logMsg string) {
	println(logMsg)
}

func init() {
	g_logger = new(consoleLog)
}

func SetLogger(logger Logger) {
	g_logger = logger
}
