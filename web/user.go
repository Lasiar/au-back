package web

import (
	"fmt"
	"github/Lasiar/au-back/model/auth"
	"github/Lasiar/au-back/web/context"
	"net/http"

	"golang.org/x/crypto/bcrypt"
)

// FullUser структура для редактирования пользователя
type FullUser struct {
	ID       int    `json:"id"`
	Login    string `json:"login"`
	Name     string `json:"name"`
	PermID   int    `json:"perm_id"`
	HashPass string `json:"hash_pass"`
	Pass     string `json:"password"`
}

// User структура для пользователя
type User struct {
	ID       int    `json:"id"`
	Login    string `json:"login"`
	Name     string `json:"name"`
	PermMask int    `json:"mask"`
}

// Load конвертирует данные из модуля auth
func (u *User) Load(user *auth.User) *User {
	u.ID = user.ID
	u.Login = user.Login
	u.PermMask = user.PermID
	if user.Name.Valid {
		u.Name = user.Name.String
	}
	return u
}

// Load конвертирует данные из модуля auth
func (u *FullUser) Load(user *auth.User) *FullUser {
	u.ID = user.ID
	u.Login = user.Login
	u.HashPass = user.Pass
	if user.Name.Valid {
		u.Name = user.Name.String
	}
	return u
}

// Upload конвертирует данные для модуля auth
func (u *FullUser) Upload() (*auth.User, error) {
	user := &auth.User{}
	user.ID = u.ID
	user.Login = u.Login
	user.Pass = u.HashPass
	if len(u.Pass) > 0 {
		pass, err := bcrypt.GenerateFromPassword([]byte(u.Pass), bcrypt.MinCost)
		if err != nil {
			return nil, err
		}
		user.Pass = string(pass)
	}
	if len(u.Name) > 0 {
		user.Name.Valid = true
		user.Name.String = u.Name
	}
	return user, nil
}

// SetUser устанавливает параметры пользователя
func SetUser() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		dr := &FullUser{}
		if err := ParseJSON(r, &dr); err != nil {
			context.SetError(r, err)
			return
		}
		user, err := dr.Upload()
		if err != nil {
			context.SetError(r, err)
			return
		}
		if len(user.Pass) < 1 {
			err := fmt.Errorf("пустой пароль")
			context.SetError(r, err)
			return
		}
		if len(user.Login) < 1 {
			err := fmt.Errorf("пустой логин")
			context.SetError(r, err)
			return
		}
		id, err := auth.GetAuth().ChangeUser(user)
		context.SetErrorOrResponse(r, id, err)
	})
}

// Login авторизирует пользователя
func Login() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		req := struct {
			Login string `json:"login"`
			Pass  string `json:"password"`
		}{}
		if err := ParseJSON(r, &req); err != nil {
			context.SetError(r, err)
			return
		}
		if token, err := GetToken(r); err == nil && token != "" {
			auth.GetAuth().CloseSession(token)
		}
		user, session, err := auth.GetAuth().Login(req.Login, req.Pass)
		if err != nil {
			context.SetError(r, err)
			return
		}
		SetToken(w, session.Token)
		context.SetResponse(r, (&User{}).Load(user))
	})
}
