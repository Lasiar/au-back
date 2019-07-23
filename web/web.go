package web

import (
	"log"
	"net/http"
	"os"
	"time"

	"github.com/Lasiar/au-back/base"
	"github.com/Lasiar/au-back/web/game"
	"github.com/Lasiar/au-back/web/middleware"
	"github.com/Lasiar/au-back/web/user"
)

// Run settings and run web server on specified port in config
func Run() {
	apiMux := http.NewServeMux()
	// работа с пользователем
	apiMux.Handle("/api/user/login", user.Login())
	apiMux.Handle("/api/user/registration", user.RegistrationUser())
	apiMux.Handle("/api/user/permissions", user.GetPermissions())
	apiMux.Handle("/api/user/logout", user.Logout())
	apiMux.Handle("/api/user/set", middleware.Permission("admin_full", user.SetUser()))
	apiMux.Handle("/api/users", middleware.Permission("admin_full", user.GetUsers()))
	apiMux.Handle("/api/user", middleware.Permission("user", user.GetUser()))
	// работа с игрой
	apiMux.Handle("/api/game/new", middleware.Permission("user", game.CreateSession()))
	apiMux.Handle("/api/game/sessions", middleware.Permission("user", game.GetSessions()))
	apiMux.Handle("/api/game/guess", middleware.Permission("user", game.Guess()))
	apiMux.Handle("/api/game/history", middleware.Permission("user", game.History()))
	apiMux.Handle("/api/game/leaderboard", middleware.Permission("user", game.Leaderboard()))

	logger := log.New(os.Stdout, "[connect] ", log.Flags())
	api := middleware.CORS("POST, GET", middleware.JSONWrite(apiMux))
	webServer := &http.Server{
		Addr:           base.GetConfig().Port,
		Handler:        middleware.Logging(logger)(api),
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	if err := webServer.ListenAndServe(); err != nil {
		log.Fatalf("Ошибка запуска сервера %v", err)
	}
}
