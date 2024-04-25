package logger

import (
	"fmt"
	"os"
	"path"
	"runtime"
	"sync"
	"time"

	"github.com/zjl410091917/game-rank/interal/app"
	"github.com/zjl410091917/game-rank/interal/config"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

const (
	defaultCap = "default"
)

type Logger struct {
	lz       *zap.Logger
	lzMap    map[string]*zap.Logger
	hostName string
}

var (
	instance *Logger
	once     sync.Once
)

func GetInstance() *Logger {
	once.Do(func() {
		instance = &Logger{
			lzMap: make(map[string]*zap.Logger),
		}
	})
	return instance
}

func (l *Logger) OnInit() {
	l.lz = l.GetCategory(defaultCap)
	app.GetInstance().SetLogger(l.lz)
	hostName, err := os.Hostname()
	if err != nil {
		hostName = "unknown"
	}
	l.hostName = hostName
}

func (l *Logger) OnStart() {
}

func (l *Logger) OnStop() {}

func (l *Logger) GetCategory(name string) *zap.Logger {
	lz := l.lzMap[name]
	if lz != nil {
		return lz
	}

	var ws []zapcore.WriteSyncer
	ws = append(ws, zapcore.AddSync(os.Stdout))
	cc := config.C()
	if len(cc.Logger.FileOutput) > 0 {
		logPath := path.Join(cc.Logger.FileOutput, app.GetInstance().Name())
		_ = os.MkdirAll(logPath, os.ModePerm)

		ws = append(ws, zapcore.AddSync(&lumberjack.Logger{
			Filename:   path.Join(logPath, fmt.Sprintf("%s.log", name)),
			MaxSize:    512,
			MaxAge:     10,
			MaxBackups: 5,
			LocalTime:  true,
			Compress:   false,
		}))
	}
	lz = zap.New(
		zapcore.NewCore(
			zapcore.NewJSONEncoder(zapcore.EncoderConfig{
				MessageKey:     "msg",
				LevelKey:       "lv",
				TimeKey:        "tm",
				NameKey:        "n",
				StacktraceKey:  "stack",
				EncodeLevel:    zapcore.LowercaseLevelEncoder,
				EncodeTime:     zapcore.TimeEncoderOfLayout(time.RFC3339Nano),
				EncodeDuration: zapcore.StringDurationEncoder,
				EncodeCaller:   zapcore.ShortCallerEncoder,
				EncodeName:     zapcore.FullNameEncoder,
			}),
			zapcore.NewMultiWriteSyncer(ws...),
			ConvertToZapLevel(config.C().Logger.Level),
		),
		zap.AddStacktrace(zapcore.ErrorLevel),
	)
	l.lzMap[name] = lz
	return lz
}

func Info(msg string, fields ...zap.Field) {
	l := GetInstance()
	var mem runtime.MemStats
	runtime.ReadMemStats(&mem)

	fields = append(fields,
		zap.Uint64("HeapInuse(MB)", mem.HeapInuse/1024/1024),
		zap.Uint64("HeapIdle(MB)", mem.HeapIdle/1024/1024),
		zap.Uint64("HeapAlloc(MB)", mem.HeapAlloc/1024/1024),
		zap.Uint64("HeapReleased(MB)", mem.HeapReleased/1024/1024),
		zap.Uint64("HeapSys(MB)", mem.HeapSys/1024/1024),
		zap.String("host-name", l.hostName),
	)
	l.lz.Info(msg, fields...)
}

func Warn(msg string, fields ...zap.Field) {
	l := GetInstance()
	l.lz.Warn(msg, fields...)
}

func Debug(msg string, fields ...zap.Field) {
	l := GetInstance()
	l.lz.Debug(msg, fields...)
}

func Error(msg string, fields ...zap.Field) {
	l := GetInstance()
	l.lz.Error(msg, fields...)
}

func ErrorWithPanic(msg string, fields ...zap.Field) {
	Error(msg, fields...)
	panic(msg)
}
