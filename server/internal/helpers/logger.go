package helpers

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gorm.io/gorm/logger"
	"os"
	"time"
)

func CustomFormatter() *zap.Logger {
	// Define a custom encoder configuration
	encoderConfig := zapcore.EncoderConfig{
		TimeKey:        "time",
		LevelKey:       "level",
		NameKey:        "logger",
		CallerKey:      "caller",
		MessageKey:     "msg",
		StacktraceKey:  "stacktrace",
		LineEnding:     "\n\n",
		EncodeLevel:    customLevelEncoder, // Custom log level format
		EncodeTime:     customTimeEncoder,  // Custom time format
		EncodeDuration: zapcore.StringDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder, // Short file name + line number
	}

	consoleEncoder := zapcore.NewConsoleEncoder(encoderConfig)
	core := zapcore.NewCore(consoleEncoder, zapcore.Lock(zapcore.AddSync(os.Stdout)), zap.DebugLevel)
	return zap.New(core, zap.AddCaller(), zap.AddStacktrace(zap.ErrorLevel))
}

func customLevelEncoder(level zapcore.Level, enc zapcore.PrimitiveArrayEncoder) {
	color := logger.Green
	switch level {
	case zapcore.DebugLevel:
		color = logger.Blue
	case zapcore.InfoLevel:
		color = logger.Green
	case zapcore.WarnLevel:
		color = logger.Yellow
	case zapcore.ErrorLevel:
		color = logger.Red
	case zapcore.DPanicLevel:
		color = logger.MagentaBold
	case zapcore.PanicLevel:
		color = logger.MagentaBold
	default:
		color = logger.Cyan
	}
	enc.AppendString(color + "[" + level.CapitalString() + "]" + logger.Reset)
}

func customTimeEncoder(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
	enc.AppendString(logger.YellowBold + "[" + t.Format("15:04:05.000") + "]" + logger.Reset)
}
