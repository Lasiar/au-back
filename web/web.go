package web

import (
	"database/sql"
	"encoding/json"
	"errors"
	"github/Lasiar/au-back/base"
	"github/Lasiar/au-back/model/auth"
	"github/Lasiar/au-back/web/context"
	"log"
	"net/http"
	"os"
	"time"
)

var (
	ErrForbidden = errors.New("not forbidden")
)

// Run settings and run web server on specified port in config
func Run() {
	apiMux := http.NewServeMux()
	apiMux.Handle("/api/set-user", CORSMiddleware("POST", SetUser()))
	apiMux.Handle("/api/login", CORSMiddleware("POST", Login()))
	logger := log.New(os.Stdout, "[connect] ", log.Flags())

	api := JSONWriteHandler(apiMux)

	webServer := &http.Server{
		Addr:           base.GetConfig().Port,
		Handler:        middlewareLogging(logger)(api),
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	if err := webServer.ListenAndServe(); err != nil {
		log.Fatalf("Ошибка запуска сервера %v", err)
	}
}

func middlewareLogging(logger *log.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()
			defer func() {
				logger.Println(r.Method, r.URL.Path, r.RemoteAddr, r.UserAgent(), "time: ", time.Since(start))
			}()
			next.ServeHTTP(w, r)
		})
	}
}

// JSONWriteHandler хандлер для ответа в виде json
func JSONWriteHandler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		if next != nil {
			next.ServeHTTP(w, r)
		}

		if err := r.Context().Err(); err != nil {
			log.Println(err)
			switch {
			case err == sql.ErrNoRows:
				http.Error(w, "Нет данных по данному запросы", http.StatusNotFound)
			default:
				http.Error(w, "error", http.StatusInternalServerError)
			}
			return
		}

		w.WriteHeader(http.StatusOK)

		data := r.Context().Value(context.ResponseDataKey)
		if data == nil {
			return
		}

		if err := json.NewEncoder(w).Encode(data); err != nil {
			log.Println(err)
		}
	})
}

func CORSMiddleware(method string, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Access-Control-Allow-Origin", "http://localhost:8081, chrome-extension://aejoelaoggembcahagimdiliamlcdmfmchrome-extension://aejoelaoggembcahagimdiliamlcdmfm")
		w.Header().Add("Access-Control-Allow-Headers", "*")
		w.Header().Add("Access-Control-Allow-Credentials", "true")
		w.Header().Add("Access-Control-Allow-Methods", "OPTIONS, "+method)
		w.Header().Add("Access-Control-Allow-Headers", "Content-Type, *")
		if r.Method == http.MethodOptions {
			w.Header().Add("Allow", "OPTIONS, "+method)
			return
		}
		next.ServeHTTP(w, r)
	})
}

// PermissionHandler проверяет наличие необходимых привелегий
func PermissionHandler(permName string, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if ok, err := HasPerm(r, permName); err != nil || !ok {
			context.SetError(r, ErrForbidden)
			return
		}
		if next != nil {
			next.ServeHTTP(w, r)
		}
	})
}

// GetToken возврашает токен из веб запроса
func GetToken(r *http.Request) (string, error) {
	cookie, err := r.Cookie("sms_manager")
	if err != nil {
		return "", err
	}
	return cookie.Value, nil
}

// ParseJSON читает из тела запроса переданную структуру
func ParseJSON(r *http.Request, data interface{}) error {
	dec := json.NewDecoder(r.Body)
	err := dec.Decode(&data)
	return err
}

// SetToken устанавливает токен в веб запрос
func SetToken(w http.ResponseWriter, token string) {
	http.SetCookie(w, &http.Cookie{
		Name:  "game",
		Value: token,
		Path:  "/",
	})
}

// HasPerm возвращает истину если у пользователя есть необходимая привелегия
func HasPerm(r *http.Request, code string) (bool, error) {
	token, err := GetToken(r)
	if err != nil {
		return false, err
	}
	user, _, err := auth.GetAuth().GetSession(token)
	if err != nil {
		return false, err
	}
	ok, err := auth.GetAuth().HasPerm(user, code)
	if err != nil {
		return false, err
	}
	return ok, nil
}
