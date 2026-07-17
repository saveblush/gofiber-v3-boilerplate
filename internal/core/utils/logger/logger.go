package logger

import (
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var Log *zap.SugaredLogger

// InitLogger init logger
func InitLogger() {
	level := zapcore.InfoLevel
	lowPriority := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		return lvl >= level
	})

	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder

	consoleEncoder := zapcore.NewJSONEncoder(encoderConfig)
	core := zapcore.NewCore(
		consoleEncoder,
		zapcore.AddSync(os.Stdout),
		lowPriority,
	)

	logger := zap.New(core, zap.AddCaller())

	// ไม่จำเป็นต้องประกาศ global เพราะทุก pkg เรียกใช้ lib นี้
	// zap.ReplaceGlobals(logger)
	// Log = zap.S()

	Log = logger.Sugar()
}
