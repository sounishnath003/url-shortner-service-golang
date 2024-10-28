package core

import (
	"log/slog"

	"github.com/sounishnath003/url-shortner-service-golang/cmd/utils"
)

type Core struct {
	Port      int
	Version   string
	JwtSecret string
	Lo        *slog.Logger
}

// InitCore helps to initialize all the necessary configuration
// to run the backend services.
func InitCore() *Core {
	return &Core{
		Version:   "0.0.1",
		Port:      utils.GetEnv("PORT", 3000).(int),
		JwtSecret: utils.GetEnv("JWT_SECRET", "secret@1234#!").(string),
		Lo:        slog.Default(),
	}
}
