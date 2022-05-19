package loggerInit

import (
	"testing"

	"go.uber.org/zap"
)

func Test_log(t *testing.T) {
	zap.L().Info("Info")
	zap.L().Debug("Debug")
	zap.L().Warn("Warn")
	zap.L().Error("Error")
}
