package auth

import (
	"crypto/sha256"
	"fmt"
	"time"

	"golang.org/x/crypto/bcrypt"
)

// Auth структура для работы с аутентификацией
type Auth struct {
	db *database
}

func hash(s string) string { return fmt.Sprintf("%x", sha256.Sum256([]byte(s))) }

// CloseSession завершает сессию, игнорирует ошибки
func (a *Auth) CloseSession(token string) { a.db.deleteSessionByToken(token) }

// NewSession создает новую сессию для пользователя
func (a *Auth) NewSession(user *User) (*Session, error) {
	session := &Session{
		UserID: user.ID,
		Token:  hash(user.Login + time.Now().String()),
	}
	session.LastUpdate.Time = time.Now()
	session.LastUpdate.Valid = true
	err := a.db.insertSession(session)
	if err != nil {
		return nil, err
	}
	return session, nil
}

// Login проверяет пользователя и в случае успеха создает сессию
func (a *Auth) Login(login, pass string) (*User, *Session, error) {
	user, err := a.db.selectUserByLogin(login)
	if err != nil {
		return nil, nil, errWrongPassword
	}
	err = bcrypt.CompareHashAndPassword([]byte(user.Pass), []byte(pass))
	if err != nil {
		return nil, nil, errWrongPassword
	}
	session, err := a.NewSession(user)
	if err != nil {
		return nil, nil, err
	}
	return user, session, nil
}

// GetSession возвращает текущую сессию пользователя или ошибку если сессии не существует
func (a *Auth) GetSession(token string) (*User, *Session, error) {
	session, err := a.db.selectSessionByToken(token)
	if err != nil {
		return nil, nil, errSession
	}
	a.db.updateSessionTime(token, time.Now())
	user, err := a.db.selectUserByID(session.UserID)
	if err != nil {
		return nil, nil, errSession
	}
	return user, session, nil
}

// HasPerm проверяет есть ли у пользователя привелегия
func (a *Auth) HasPerm(user *User, code string) (bool, error) {
	perm, err := a.db.selectPermissionByCode(code)
	if err != nil {
		return false, err
	}
	if perm.ID > user.PermID {
		return false, nil
	}
	return true, nil
}

// GetPermissions возвращает список прав доступа
func (a *Auth) GetPermissions() ([]*Permission, error) {
	return a.db.selectPermissions()
}
