package logging

import (
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type SugarLogger interface {
	Infow(msg string, kv ...interface{})
	Errorw(msg string, kv ...interface{})
	Fatalw(msg string, kv ...interface{})
}

var _ SugarLogger = (*zap.SugaredLogger)(nil)

func NewZapSugarLogger() *zap.SugaredLogger {
	encoderConfig := encoderConfig()

	// logFile, err := os.OpenFile(
	// 	fmt.Sprintf("logs/ap-%d.log", time.Now().Unix()),
	// 	os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644,
	// )

	// if err != nil {
	// 	log.Printf("failed to create log file: %v", err)
	// }

	// fileEncoder := zapcore.NewJSONEncoder(encoderConfig)
	consoleEncoder := zapcore.NewConsoleEncoder(encoderConfig)

	core := zapcore.NewTee(
		// zapcore.NewCore(fileEncoder, zapcore.AddSync(logFile), zap.DebugLevel),
		zapcore.NewCore(consoleEncoder, zapcore.AddSync(os.Stdout), zap.DebugLevel),
	)

	logger := zap.New(core, zap.AddCaller(), zap.AddCallerSkip(1), zap.AddStacktrace(zapcore.ErrorLevel))
	sugaredLogger := logger.Sugar()
	return sugaredLogger

}

func encoderConfig() zapcore.EncoderConfig {
	cfg := zap.NewProductionEncoderConfig()
	cfg.EncodeLevel = zapcore.CapitalColorLevelEncoder
	cfg.EncodeTime = zapcore.ISO8601TimeEncoder
	return cfg
}
