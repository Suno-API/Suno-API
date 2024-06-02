package common

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sync"
	"time"

	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var Logger *zap.SugaredLogger
var LoggerZap *zap.Logger
var setupLogLock sync.Mutex

func SetupLogger() {
	if Logger != nil {
		return
	}
	ok := setupLogLock.TryLock()
	if !ok {
		log.Println("setup log is already working")
		return
	}
	defer func() {
		setupLogLock.Unlock()
	}()
	levelEnabler := zapcore.DebugLevel
	syncers := []zapcore.WriteSyncer{zapcore.AddSync(os.Stdout)}

	if LogDir != "" {
		var err error
		LogDir, err = filepath.Abs(LogDir)
		if err != nil {
			log.Fatal(err)
		}
		if _, err := os.Stat(LogDir); os.IsNotExist(err) {
			err = os.Mkdir(LogDir, 0777)
			if err != nil {
				log.Fatal(err)
			}
		}

		if RotateLogs {
			fd, err := rotatelogs.New(
				filepath.Join(LogDir, "%Y%m%d", "%H:%M.log"),
				rotatelogs.WithRotationTime(time.Hour),
				rotatelogs.WithMaxAge(time.Hour*24*100),
			)
			if err != nil {
				log.Fatal("failed to open rotateLogs")
			}
			syncers = append(syncers, zapcore.AddSync(fd))
		} else {
			logPath := filepath.Join(LogDir, fmt.Sprintf("any2api-%s.log", time.Now().Format("20060102")))
			fd, err := os.OpenFile(logPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
			if err != nil {
				log.Fatal("failed to open log file")
			}
			syncers = append(syncers, zapcore.AddSync(fd))
		}
	}
	enc := zapcore.NewJSONEncoder(zapcore.EncoderConfig{
		TimeKey:   "_time", // Modified
		LevelKey:  "level",
		NameKey:   "logger",
		CallerKey: "caller",
		//FunctionKey:    zapcore.OmitKey,
		MessageKey:     "msg",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.LowercaseLevelEncoder,
		EncodeTime:     zapcore.RFC3339NanoTimeEncoder, // Modified
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	})

	loggerCore := zapcore.NewCore(enc,
		zapcore.NewMultiWriteSyncer(syncers...),

		levelEnabler)
	LoggerZap = zap.New(loggerCore,
		zap.AddStacktrace(
			zap.NewAtomicLevelAt(zapcore.ErrorLevel)),
		zap.AddCaller(),
		zap.AddCallerSkip(0),
	)
	Logger = LoggerZap.Sugar()
	//gin.DefaultWriter = Logger.Writer()
	//gin.DefaultErrorWriter = io.MultiWriter(os.Stderr, logsWriter)
}
