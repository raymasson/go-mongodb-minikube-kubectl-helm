package logger

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

//Log ...
var Log *zap.Logger

//New : instantiate a new logger
func New() {
	conf := zap.NewProductionConfig()
	conf.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	conf.EncoderConfig.MessageKey = "message"
	conf.EncoderConfig.TimeKey = "timestamp"

	Log, _ = conf.Build()
}
