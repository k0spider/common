package psql

import (
	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
	_ "github.com/lib/pq"
	"github.com/sirupsen/logrus"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"time"
)

type GormConfig struct {
	Driver      string        `yaml:"driver"`
	Dsn         []string      `yaml:"dsn"`
	IdleTimeout time.Duration `yaml:"idleTimeout"`
	MaxOpens    int           `yaml:"maxOpens"`
	MaxIdles    int           `yaml:"maxIdles"`
	Debug       bool          `yaml:"debug"`
	LogDir      string        `yaml:"logDir"`
	LogMaxDay   int           `yaml:"logMaxDay"`
}

func NewGorm(c *GormConfig) *gorm.DB {
	// 创建日志目录
	logPath := "./logs/gorm"
	if c.LogDir != "" {
		logPath = c.LogDir
	}
	if c.LogMaxDay == 0 {
		c.LogMaxDay = 30
	}
	writer, _ := rotatelogs.New(
		logPath+"/%Y%m%d.log",
		rotatelogs.WithMaxAge(time.Duration(c.LogMaxDay)*24*time.Hour), // 保留 30 天的日志
		rotatelogs.WithRotationTime(24*time.Hour),                      // 每天分割一次日志
	)
	log := logrus.New()
	log.SetOutput(writer)
	log.SetFormatter(&RawLogFormatter{})
	newLogger := New(
		NewGormLogger(log), // io writer
		logger.Config{
			SlowThreshold:             time.Second,   // Slow SQL threshold
			LogLevel:                  logger.Silent, // Log level
			IgnoreRecordNotFoundError: false,         // Ignore ErrRecordNotFound error for logger
			ParameterizedQueries:      false,
			Colorful:                  false,
		},
	)
	db, err := gorm.Open(postgres.Open(c.Dsn[0]), &gorm.Config{
		SkipDefaultTransaction: true, //禁用默认事务
		PrepareStmt:            true, //缓存预编译语句
		Logger:                 newLogger,
	})
	if err != nil {
		panic(err)
	}
	sqlDB, err := db.DB()
	if err != nil {
		panic(err)
	}
	sqlDB.SetMaxIdleConns(c.MaxIdles)
	sqlDB.SetMaxOpenConns(c.MaxOpens)
	sqlDB.SetConnMaxLifetime(c.IdleTimeout)
	if c.Debug {
		db = db.Debug()
	}
	if err = sqlDB.Ping(); err != nil {
		panic(err)
	}

	return db
}

// NewGormLogger 创建一个自定义的 GORM 日志记录器
func NewGormLogger(log *logrus.Logger) *GormLogger {
	return &GormLogger{log}
}

// GormLogger 是一个自定义的 GORM 日志记录器
type GormLogger struct {
	Log *logrus.Logger
}

// Print 实现 GORM 的日志记录接口
func (l *GormLogger) Printf(format string, v ...any) {
	l.Log.Printf(format, v...)

}

// RawLogFormatter 是自定义的 logrus.Formatter，用于输出原始日志消息
type RawLogFormatter struct{}

// Format 格式化日志消息，这里直接返回原始消息
func (f *RawLogFormatter) Format(entry *logrus.Entry) ([]byte, error) {
	nowTime := time.Now().Format(time.DateTime)
	return []byte(nowTime + "  " + entry.Message + "\n\n"), nil
}
