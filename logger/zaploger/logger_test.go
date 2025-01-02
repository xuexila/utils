package zaploger

import (
	"gorm.io/gorm/logger"
	"testing"

	"go.uber.org/zap/zapcore"
)

func TestConvertLogLevel(t *testing.T) {
	tests := []struct {
		name     string
		level    logger.LogLevel
		expected zapcore.LevelEnabler
	}{
		{"Silent", logger.Silent, zapcore.PanicLevel},
		{"Error", logger.Error, zapcore.ErrorLevel},
		{"Warn", logger.Warn, zapcore.WarnLevel},
		{"Info", logger.Info, zapcore.DebugLevel},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := convertLogLevel(tt.level)
			if got != tt.expected {
				t.Errorf("convertLogLevel(%v) = %v, want %v", tt.level, got, tt.expected)
			}
		})
	}
}
