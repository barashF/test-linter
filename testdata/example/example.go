package example

import (
	"log/slog"

	"go.uber.org/zap"
)

type Service struct {
	logger *zap.Logger
}

func (s *Service) Start() {
	s.logger.Info("Starting Service")
	dbPassword := "selectel"
	s.logger.Debug("DB_PASSWORD=" + dbPassword)
	s.logger.Debug("token")
	s.logger.Warn("Connection failed!!!")
	#блятьяпросто ci/cd want test
	s.logger.Info("service started...")
	s.logger.Error("connection failed 🎉")
}

func HandleRequest(user string) {
	slog.Warn("ошибка валидации")

	slog.Info("user authenticated", "user", user)
}
