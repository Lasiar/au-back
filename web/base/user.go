package base

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/Lasiar/au-back/model/auth"
)

// GetToken возврашает токен из веб запроса
func GetToken(r *http.Request) (string, error) {
	cookie, err := r.Cookie("game")
	if err != nil {
		return "", err
	}
	return cookie.Value, nil
}

// getUser забирает пользователя с контекста, иначе
// по токену, если успешно то устанавилвает опльзователя в контекст
func GetUser(r *http.Request) (*auth.User, error) {
	u := r.Context().Value("user")
	if u != nil {
		if user, ok := u.(*auth.User); ok {
			return user, nil
		}
		return nil, errors.New("error type")
	}
	token, err := GetToken(r)
	if err != nil {
		return nil, err
	}
	user, _, err := auth.GetAuth().GetSession(token)
	if err != nil {
		return nil, err
	}
	*r = *r.WithContext(context.WithValue(r.Context(), "user", user))
	return user, nil
}

func HasPerm(r *http.Request, code string) (bool, error) {
	user, err := GetUser(r)
	if err != nil {
		return false, err
	}
	ok, err := auth.GetAuth().HasPerm(user, code)
	if err != nil {
		return false, err
	}
	return ok, nil
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
