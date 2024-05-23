package log

import (
	"context"
	"github.com/k0spider/common/utils"
	"time"

	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
	"github.com/rifflock/lfshook"
	"github.com/sirupsen/logrus"
)

type Logs struct {
	Level         string        `yaml:"level"`
	Path          string        `yaml:"path"`
	FileSuffix    string        `yaml:"fileSuffix"`
	MaxAgeHour    time.Duration `yaml:"maxAgeHour"`
	RotationCount uint          `yaml:"rotationCount"`
}

var DefaultLogger *logrus.Logger

func InitLogger(logs *Logs) {
	DefaultLogger = logrus.New()
	level, _ := logrus.ParseLevel(logs.Level)
	DefaultLogger.SetLevel(level)
	DefaultLogger.SetReportCaller(false)
	DefaultLogger.SetFormatter(&logrus.JSONFormatter{TimestampFormat: "2006-01-02 15:04:05", DisableTimestamp: false, PrettyPrint: false})
	writer, _ := rotatelogs.New(logs.Path+"%Y%m%d"+logs.FileSuffix,
		rotatelogs.WithRotationTime(24*time.Hour),
		rotatelogs.WithRotationCount(logs.RotationCount),
		rotatelogs.WithRotationTime(time.Hour*logs.MaxAgeHour),
	)
	DefaultLogger.AddHook(lfshook.NewHook(writer, &logrus.JSONFormatter{TimestampFormat: "2006-01-02 15:04:05"}))
}
func WithContext(ctx context.Context) *logrus.Entry {
	if ctx.Value(utils.RequestID) != nil {
		return DefaultLogger.WithContext(ctx).WithField(utils.RequestID, ctx.Value(utils.RequestID))
	}
	return DefaultLogger.WithContext(ctx)
}
