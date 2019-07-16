package middleware

import (
	"database/sql"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"time"

	web "github.com/Lasiar/au-back/web/base"
	"github.com/Lasiar/au-back/web/context"
)

// JSONWrite хандлер для ответа в виде json
func JSONWrite(next http.Handler) http.Handler {
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

func CORS(method string, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Access-Control-Allow-Origin", "http://localhost:8080")
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

// Permission проверяет наличие необходимых привелегий
func Permission(permName string, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if ok, err := web.HasPerm(r, permName); err != nil || !ok {
			context.SetError(r, errors.New("forrbiden"))
			return
		}

		if next != nil {
			next.ServeHTTP(w, r)
		}
	})
}

func Logging(logger *log.Logger) func(http.Handler) http.Handler {
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
