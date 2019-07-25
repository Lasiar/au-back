package user

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/Lasiar/au-back/model/auth"
	web "github.com/Lasiar/au-back/web/base"
	"github.com/Lasiar/au-back/web/context"

	"golang.org/x/crypto/bcrypt"
)

// FullUser структура для редактирования пользователя
type FullUser struct {
	ID       int    `json:"id"`
	Login    string `json:"login"`
	Name     string `json:"name"`
	PermID   int    `json:"perm_mask"`
	HashPass string `json:"hash_pass"`
	Pass     string `json:"password"`
}

// User структура для пользователя
type User struct {
	ID       int    `json:"id"`
	Login    string `json:"login"`
	Name     string `json:"name"`
	PermMask int    `json:"perm_mask"`
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
	u.PermID = user.PermID
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
	user.PermID = u.PermID
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
		if err := web.ParseJSON(r, &dr); err != nil {
			context.SetError(r, err)
			return
		}
		user, err := dr.Upload()
		if err != nil {
			context.SetError(r, err)
			return
		}
		if user.Login == "" {
			context.SetError(r, fmt.Errorf("%v:%v", web.ErrBadRequest, "empty login"))
			return
		}
		id, err := auth.GetAuth().ChangeUser(user)
		if err != nil && strings.Contains(err.Error(), "users_id_permission_fkey") {
			context.SetError(r, fmt.Errorf("%v: %v", web.ErrBadRequest, "perm mask does not exist`"))
			return
		}
		context.SetErrorOrResponse(r, id, err)
	})
}

// GetUser отдает пользователя по токену, если нет, то 403
func GetUser() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token, err := web.GetToken(r)
		if err != nil {
			context.SetError(r, err)
			return
		}
		user, _, err := auth.GetAuth().GetSession(token)
		if err != nil {
			context.SetError(r, err)
		}
		context.SetResponse(r, (&User{}).Load(user))
	})
}

// SetUser устанавливает параметры пользователя
func RegistrationUser() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		dr := struct {
			FullUser
			Captcha string `json:"captcha"`
		}{}
		if err := web.ParseJSON(r, &dr); err != nil {
			context.SetError(r, err)
			return
		}
		user, err := dr.Upload()
		if err != nil {
			context.SetError(r, err)
			return
		}
		if user.Pass == "" || user.Login == "" {
			context.SetError(r, fmt.Errorf("%v:%v", web.ErrBadRequest, "empty login or password"))
			return
		}
		id, err := auth.GetAuth().AddUser(user)
		if err != nil && strings.Contains(err.Error(), "users_login_key") {
			context.SetError(r, fmt.Errorf("%v: %v", web.ErrBadRequest, "Такой ник уже занят"))
			return
		}
		context.SetErrorOrResponse(r, id, err)
	})
}

// Login авторизирует пользователя
func Login() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		req := &struct {
			Login string `json:"login"`
			Pass  string `json:"password"`
		}{}
		if err := web.ParseJSON(r, req); err != nil {
			context.SetError(r, err)
			return
		}
		if req.Pass == "" || req.Login == "" {
			context.SetError(r, fmt.Errorf("%v:%v", web.ErrBadRequest, "empty login"))
			return
		}
		if token, err := web.GetToken(r); err == nil && token != "" {
			auth.GetAuth().CloseSession(token)
		}
		_, session, err := auth.GetAuth().Login(req.Login, req.Pass)
		if err != nil {
			context.SetError(r, err)
			return
		}
		web.SetToken(w, session.Token)
		context.SetResponse(r, struct {
			Token string `json:"token"`
		}{Token: session.Token})
	})
}

func GetPermissions() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		type permission struct {
			Level    int    `json:"level"`
			CodeName string `json:"code_name"`
		}
		perms, err := auth.GetAuth().GetPermissions()
		if err != nil {
			context.SetError(r, err)
			return
		}
		var resp []permission
		for _, p := range perms {
			resp = append(resp, permission{Level: p.ID, CodeName: p.Code})
		}
		context.SetResponse(r, resp)
	})
}

func Logout() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token, err := web.GetToken(r)
		if err != nil {
			context.SetError(r, err)
			return
		}
		auth.GetAuth().CloseSession(token)
		context.SetResponse(r, struct{}{})
	})
}

// GetUsers возвращает список пользователей
func GetUsers() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		dUsers, err := auth.GetAuth().GetUsers()
		if err != nil {
			context.SetError(r, err)
			return
		}
		users := make([]*User, 0)
		for _, dUser := range dUsers {
			users = append(users, (&User{}).Load(dUser))
		}
		context.SetResponse(r, &users)
	})
}
