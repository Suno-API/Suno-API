package gormlogger

import (
	"context"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/utils"
	"runtime"
	"strconv"
	"strings"
	"time"
)

type GormV2LogLevel struct {
	GormV2 GormV2LoggerConfig `json:"gorm_v2"`
}

type GormV2LoggerConfig struct {
	SlowThreshold time.Duration   `json:"slow_threshold"`
	LogLevel      logger.LogLevel `json:"log_level"`
	Colorful      bool            `json:"colorful"`
}

func NewGormV2Logger(zaplogger *zap.SugaredLogger) (logger.Interface, error) {
	gormV2LogLevel := GormV2LogLevel{
		GormV2: GormV2LoggerConfig{
			SlowThreshold: time.Second,
			LogLevel:      logger.Info,
			Colorful:      true,
		},
	}

	switch zaplogger.Level() {
	case zapcore.ErrorLevel:
		gormV2LogLevel.GormV2.LogLevel = logger.Error
	case zapcore.WarnLevel:
		gormV2LogLevel.GormV2.LogLevel = logger.Warn
	}

	gormLogger := New(
		zaplogger,
		logger.Config{
			SlowThreshold: gormV2LogLevel.GormV2.SlowThreshold,
			LogLevel:      gormV2LogLevel.GormV2.LogLevel,
			Colorful:      gormV2LogLevel.GormV2.Colorful,
		},
	)
	return gormLogger, nil
}

type ctxLogger struct {
	*zap.SugaredLogger

	logger.Config
	infoStr, warnStr, errStr            string
	traceStr, traceErrStr, traceWarnStr string
}

// New new ctxLogger
func New(writer *zap.SugaredLogger, config logger.Config) logger.Interface {
	var (
		infoStr      = "%s\n[info] "
		warnStr      = "%s\n[warn] "
		errStr       = "%s\n[error] "
		traceStr     = "%s\n[%v] [rows:%d] %s"
		traceWarnStr = "%s\n[%v] [rows:%d] %s"
		traceErrStr  = "%s %s\n[%v] [rows:%d] %s"
	)

	if config.Colorful {
		infoStr = logger.Green + "%s\n" + logger.Reset + logger.Green + "[info] " + logger.Reset
		warnStr = logger.Blue + "%s\n" + logger.Reset + logger.Magenta + "[warn] " + logger.Reset
		errStr = logger.Magenta + "%s\n" + logger.Reset + logger.Red + "[error] " + logger.Reset
		traceStr = logger.Green + "%s\n" + logger.Reset + logger.Yellow + "[%.3fms] " + logger.Blue + "[rows:%d]" + logger.Reset + " %s"
		traceWarnStr = logger.Green + "%s\n" + logger.Reset + logger.RedBold + "[%.3fms] " + logger.Yellow + "[rows:%d]" + logger.Magenta + " %s" + logger.Reset
		traceErrStr = logger.RedBold + "%s " + logger.MagentaBold + "%s\n" + logger.Reset + logger.Yellow + "[%.3fms] " + logger.Blue + "[rows:%d]" + logger.Reset + " %s"
	}

	return &ctxLogger{
		SugaredLogger: writer,
		Config:        config,
		infoStr:       infoStr,
		warnStr:       warnStr,
		errStr:        errStr,
		traceStr:      traceStr,
		traceWarnStr:  traceWarnStr,
		traceErrStr:   traceErrStr,
	}
}

// LogMode log mode
func (l *ctxLogger) LogMode(level logger.LogLevel) logger.Interface {
	newlogger := *l
	newlogger.LogLevel = level
	return &newlogger
}

// Info print info
func (l ctxLogger) Info(ctx context.Context, msg string, data ...interface{}) {
	if l.LogLevel >= logger.Info {
		l.SugaredLogger.Infof(l.infoStr+msg, append([]interface{}{utils.FileWithLineNum()}, data...)...)
	}
}

// Warn print warn messages
func (l ctxLogger) Warn(ctx context.Context, msg string, data ...interface{}) {
	if l.LogLevel >= logger.Info {
		l.SugaredLogger.Warnf(l.warnStr+msg, append([]interface{}{utils.FileWithLineNum()}, data...)...)
	}
}

// Error print error messages
func (l ctxLogger) Error(ctx context.Context, msg string, data ...interface{}) {
	if l.LogLevel >= logger.Info {
		l.SugaredLogger.Errorf(l.errStr+msg, append([]interface{}{utils.FileWithLineNum()}, data...)...)
	}
}

// Trace print sql message
func (l ctxLogger) Trace(ctx context.Context, begin time.Time, fc func() (string, int64), err error) {
	if l.LogLevel > 0 {
		elapsed := time.Since(begin)
		gormLogger := l.SugaredLogger.With("root_caller", getRootCaller())
		switch {
		case err != nil && l.LogLevel >= logger.Info:
			sql, rows := fc()
			gormLogger.Debugw(sql, "rows affected", rows, "err", err)
		case elapsed > l.SlowThreshold && l.SlowThreshold != 0 && l.LogLevel >= logger.Info:
			sql, rows := fc()
			gormLogger.Debugw(sql, "rows affected", rows, "exceed threshold", true, "threshold", l.SlowThreshold, "elapsed", elapsed.String())
		case l.LogLevel >= logger.Info:
			sql, rows := fc()
			gormLogger.Debugw(sql, "rows affected", rows, "elapsed", elapsed.String())
		}
	}
}

func getRootCaller() string {
	for i := 2; i < 15; i++ {
		_, file, line, ok := runtime.Caller(i)

		if ok && !strings.Contains(file, "gorm.io/gorm") {
			_, file, line, ok = runtime.Caller(i + 1)
			return file + ":" + strconv.FormatInt(int64(line), 10)
		}
	}
	return ""
}
