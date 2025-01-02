package zaploger

import (
	"context"
	"github.com/helays/utils/tools"
	"github.com/natefinch/lumberjack"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gorm.io/gorm/logger"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// LogConfig defines the configuration for each log level.
type LogConfig struct {
	FileName   string `json:"file_name" yaml:"file_name" int:"file_name"`       // 日志文件名称
	FilePath   string `json:"file_path" yaml:"file_path" int:"file_path"`       // 日志路径
	MaxSize    int    `json:"max_size" yaml:"max_size" int:"max_size"`          // 文件最大容量 单位MB
	MaxBackups int    `json:"max_backups" yaml:"max_backups" int:"max_backups"` // 文件最多数量
	MaxAge     int    `json:"max_age" yaml:"max_age" int:"max_age"`             // 日志保留最大时长 单位天，0为不限制
	Compress   bool   `json:"compress" yaml:"compress" int:"compress"`          // 是否压缩
	ToStdout   bool   `json:"to_stdout" yaml:"to_stdout" ini:"to_stdout"`       // 是否输出到标准输出
}

// Config holds the application configuration.
type Config struct {
	LogLevel        string               `json:"log_level" yaml:"log_level" ini:"log_level"`
	LogLevelConfigs map[string]LogConfig `json:"log_level_configs" yaml:"log_level_configs" ini:"log_level_configs"` // key is log level, value is its configuration
}

// Logger implements the gorm.Logger interface.
type Logger struct {
	logger *zap.Logger
	level  zapcore.LevelEnabler
}

// customTimeEncoder formats time as "2006-01-02 15:04:05".
func customTimeEncoder(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
	enc.AppendString(t.Format(time.DateTime))
}

// New creates a new instance of Logger with given configuration.
func New(cfg *Config) (*Logger, error) {
	defalutLevel := zapcore.DebugLevel
	if cfg.LogLevel != "" {
		level, err := zapcore.ParseLevel(cfg.LogLevel)
		if err != nil {
			return nil, err
		}
		defalutLevel = level
	}
	var cores []zapcore.Core
	encoder := zapcore.NewJSONEncoder(zapcore.EncoderConfig{
		TimeKey:        "time",
		LevelKey:       "level",
		NameKey:        "logger",
		MessageKey:     "msg",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.LowercaseLevelEncoder,
		EncodeTime:     customTimeEncoder, // 使用自定义的时间编码器
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	})
	for levelStrs, config := range cfg.LogLevelConfigs {
		for _, levelStr := range strings.Split(levelStrs, ",") {
			level, err := zapcore.ParseLevel(levelStr)
			if err != nil {
				return nil, err
			}
			writers := make([]zapcore.WriteSyncer, 0)
			if config.FilePath != "" {
				fame := filepath.Join(tools.Fileabs(config.FilePath), levelStr+"_"+tools.Ternary(config.FileName == "", "log", config.FileName))
				lumberJackLogger := &lumberjack.Logger{
					Filename:   fame,
					MaxSize:    config.MaxSize,
					MaxAge:     config.MaxAge,
					MaxBackups: config.MaxBackups,
					LocalTime:  true,
					Compress:   config.Compress,
				}
				writers = append(writers, zapcore.AddSync(lumberJackLogger))
			}
			if config.ToStdout {
				writers = append(writers, zapcore.AddSync(os.Stdout))
			}
			core := zapcore.NewCore(encoder, zapcore.NewMultiWriteSyncer(writers...), level)
			cores = append(cores, core)
		}
	}
	combinedCore := zapcore.NewTee(cores...)
	return &Logger{
		logger: zap.New(combinedCore, zap.AddCaller(), zap.AddStacktrace(zap.ErrorLevel)),
		level:  defalutLevel, // Default to DebugLevel
	}, nil
}

func (l *Logger) LogMode(level logger.LogLevel) logger.Interface {
	enabler := convertLogLevel(level)
	return &Logger{
		logger: l.logger,
		level:  enabler,
	}
}

// ConvertLogLevel converts GORM log levels to Zap LevelEnabler.
func convertLogLevel(level logger.LogLevel) zapcore.LevelEnabler {
	switch level {
	case logger.Silent:
		return zapcore.PanicLevel
	case logger.Error:
		return zapcore.ErrorLevel
	case logger.Warn:
		return zapcore.WarnLevel
	default:
		return zapcore.DebugLevel // Use DebugLevel for Info and other levels
	}
}

func (l *Logger) Info(ctx context.Context, msg string, data ...interface{}) {
	if l.level.Enabled(zapcore.InfoLevel) {
		l.logger.Info(msg, zap.Any("data", data))
	}
}

func (l *Logger) Warn(ctx context.Context, msg string, data ...interface{}) {
	if l.level.Enabled(zapcore.WarnLevel) {
		l.logger.Warn(msg, zap.Any("data", data))
	}
}

func (l *Logger) Error(ctx context.Context, msg string, data ...interface{}) {
	if l.level.Enabled(zapcore.ErrorLevel) {
		l.logger.Error(msg, zap.Any("data", data))
	}
}

func (l *Logger) Trace(ctx context.Context, begin time.Time, fc func() (string, int64), err error) {
	if l.level.Enabled(zapcore.DebugLevel) {
		sql, rows := fc()
		l.logger.Debug(
			"SQL Trace",
			zap.String("sql", sql),
			zap.Int64("rows", rows),
			zap.Duration("took", time.Since(begin)),
			zap.Error(err),
		)
	}
}
