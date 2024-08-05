package logger

import (
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func ConfigZap() *zap.SugaredLogger {

	cfg := zap.Config{
		Encoding:    "json",                              //encode kiểu json hoặc console
		Level:       zap.NewAtomicLevelAt(zap.InfoLevel), //chọn InfoLevel có thể log ở cả 3 level
		OutputPaths: []string{"stderr"},

		EncoderConfig: zapcore.EncoderConfig{ //Cấu hình logging, sẽ không có stacktracekey
			MessageKey:   "message",
			TimeKey:      "time",
			LevelKey:     "level",
			CallerKey:    "caller",
			EncodeCaller: zapcore.FullCallerEncoder, //Lấy dòng code bắt đầu log
			EncodeLevel:  CustomLevelEncoder,        //Format cách hiển thị level log
			EncodeTime:   SyslogTimeEncoder,         //Format hiển thị thời điểm log
		},
	}

	logger, _ := cfg.Build() //Build ra Logger
	return logger.Sugar()    //Trả về logger hoặc Sugaredlogger, ở đây ta chọn trả về Logger
}

func SyslogTimeEncoder(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
	enc.AppendString(t.Format("2006-01-02 15:04:05"))
}

func CustomLevelEncoder(level zapcore.Level, enc zapcore.PrimitiveArrayEncoder) {
	enc.AppendString("[" + level.CapitalString() + "]")
}
