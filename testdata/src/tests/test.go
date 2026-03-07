package tests

import (
	"log/slog"

	"go.uber.org/zap"
)

func TestLowerCase() {
	slog.Info("Starting server")       // want `lowercase`
	slog.Error("Failed to connect")    // want `lowercase`
	zap.S().Warn("Connection timeout") // want `lowercase`
	zap.S().Debug("Retry attempt")     // want `lowercase`

	slog.Info("starting Server")

	slog.Info("starting server")
	slog.Error("failed to connect")
	zap.S().Warn("connection timeout")
	zap.S().Debug("retry attempt")

	msg := "Dynamic message"
	slog.Info(msg)
	zap.S().Error(msg)
}

func TestEnglishOnly() {
	slog.Info("сервер запущен")            // want `English`
	slog.Error("ошибка подключения")       // want `English`
	zap.S().Warn("предупреждение системы") // want `English`

	slog.Info("server started")
	slog.Error("connection failed")
	zap.S().Warn("system warning")

	slog.Info("version 2.0 released")
}

func TestSpecialSymbols() {
	slog.Info("server started!")        // want `special character`
	zap.S().Error("connection failed!") // want `special character`

	slog.Info("deployed 🚀")               // want `emoji or symbol`
	zap.S().Debug("status: ✅ complete.d") // want `emoji or symbol`

	slog.Warn("waiting...")    // want `ellipsis`
	zap.S().Info("loading...") // want `ellipsis`

	slog.Info("version 2.0, build 123")
	slog.Debug("user: admin, role: operator")
	zap.S().Warn("check config: timeout=30s")

	slog.Debug("price: 99.99, currency: USD")
}

func TestSensitiveData() {
	slog.Info("auth: token=abc123")       // want `sensitive.*token`
	slog.Debug("config: api_key=secret")  // want `sensitive.*api_key`
	slog.Warn("password: SuperSecret123") // want `sensitive.*password`
	zap.S().Error("credential leaked")    // want `sensitive.*credential`

	secret := "my-password"
	slog.Info("token: " + secret) // want `sensitive.*token`
	slog.Info(secret + "token")
}

func TestEdgeCases() {
	slog.Info("")
	slog.Info("a")
	slog.Info("A") // want `lowercase`

	slog.Info("server запущен") // want `English`

	slog.Info("Ошибка! 🚀") // want `lowercase` `English` `special character`

	slog.Info("math: 2 + 2 = 4")
	slog.Debug("comparison: a < b")

	logger := zap.S()
	logger.Info("starting")
	logger.Error("Failed") // want `lowercase`
}

func TestAllLogMethods() {
	slog.Debug("Debug message") // want `lowercase`
	slog.Info("Info message")   // want `lowercase`
	slog.Warn("Warn message")   // want `lowercase`
	slog.Error("Error message") // want `lowercase`

	zap.S().Debug("Debug") // want `lowercase`
	zap.S().Info("Info")   // want `lowercase`
	zap.S().Warn("Warn")   // want `lowercase`
	zap.S().Error("Error") // want `lowercase`
	zap.S().Panic("Panic") // want `lowercase`
	zap.S().Fatal("Fatal") // want `lowercase`

	slog.Debug("debug ok")
	slog.Info("info ok")
	zap.S().Warn("warn ok")
	zap.S().Error("error ok")
}
