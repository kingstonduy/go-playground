package logger

import (
	"fmt"
	"strings"
	"time"

	"github.com/fatih/color"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var LoggingMessage string = "[%s] [%s] [%s] [%s] [%s] [%s] [%s] - %s"

func NewLoggingMessage(msg string) string {
	return fmt.Sprintf(LoggingMessage, "", "", "", "", "", "", "", msg)
}

func Init() *zap.Logger {

	cfg := zap.Config{
		Encoding:    "console",                           //encode kiểu json hoặc console
		Level:       zap.NewAtomicLevelAt(zap.InfoLevel), //chọn InfoLevel có thể log ở cả 3 level
		OutputPaths: []string{"stderr"},

		EncoderConfig: zapcore.EncoderConfig{ //Cấu hình logging, sẽ không có stacktracekey
			MessageKey:       "message",
			TimeKey:          "time",
			LevelKey:         "level",
			CallerKey:        "caller",
			StacktraceKey:    "stacktrace",
			EncodeCaller:     zapcore.FullCallerEncoder, //Lấy dòng code bắt đầu log
			EncodeLevel:      CustomLevelEncoder,        //Format cách hiển thị level log
			EncodeTime:       SyslogTimeEncoder,         //Format hiển thị thời điểm log
			ConsoleSeparator: " ",
		},
	}

	logger, _ := cfg.Build() //Build ra Logger
	return logger            //Trả về logger hoặc Sugaredlogger, ở đây ta chọn trả về Logger
}

func SyslogTimeEncoder(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
	enc.AppendString(t.Format("2006-01-02 15:04:05"))
}

func CustomLevelEncoder(level zapcore.Level, enc zapcore.PrimitiveArrayEncoder) {
	green := color.New(color.FgGreen).SprintFunc()
	red := color.New(color.FgRed).SprintFunc()
	yellow := color.New(color.FgYellow).SprintFunc()
	blue := color.New(color.FgBlue).SprintFunc()

	switch level.String() {
	case "info":
		enc.AppendString(green(strings.ToUpper(level.String())))
	case "error":
		enc.AppendString(red(strings.ToUpper(level.String())))
	case "warn":
		enc.AppendString(yellow(strings.ToUpper(level.String())))
	case "debug":
		enc.AppendString(blue(strings.ToUpper(level.String())))
	default:
		enc.AppendString(strings.ToUpper(level.String()))
	}
}
