package logger

import (
	"fmt"

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
