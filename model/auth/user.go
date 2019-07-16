package auth

import (
	"database/sql"
	"fmt"
)

// User структура для отображения данных пользователя
type User struct {
	ID     int
	Login  string
	Pass   string
	PermID int
	Roles  []string
	Name   sql.NullString
}

func (u User) String() string { return fmt.Sprintf("%d %s", u.ID, u.Login) }

// selectUsers возвращает список пользователей
func (db *database) selectUsers() ([]*User, error) {
	users := make([]*User, 0)
	rows, err := db.db.Query(`SELECT id, login, pass, name, id_permission FROM auth.users`)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		user := &User{}
		if err := rows.Scan(
			&user.ID,
			&user.Login,
			&user.Pass,
			&user.Name,
			&user.PermID,
		); err != nil {
			return nil, err
		} else {
			users = append(users, user)
		}
	}
	return users, nil
}

// selectUserByID возвращает пользователя по ID
func (db *database) selectUserByID(id int) (*User, error) {
	user := &User{}
	row := db.db.QueryRow(`SELECT id, login, pass, name, id_permission FROM auth.users WHERE id = $1`, id)
	err := row.Scan(
		&user.ID,
		&user.Login,
		&user.Pass,
		&user.Name,
		&user.PermID,
	)
	if err != nil {
		return nil, err
	}
	return user, nil
}

// insertUser добавляет новго пользователя
func (db *database) insertUser(user *User) (int, error) {
	row := db.db.QueryRow(`INSERT INTO auth.users (login, pass, name, id_permission) VALUES ($1, $2, $3, $4) RETURNING id`,
		user.Login,
		user.Pass,
		user.Name,
		user.PermID,
	)
	err := row.Scan(&user.ID)
	if err != nil {
		return 0, err
	}
	return user.ID, nil
}

// updateUser обновляет данные пользователя
func (db *database) updateUser(user *User) error {
	_, err := db.db.Exec(`UPDATE auth.users SET login = $1, pass = $2, name = $3, id_permission=$4 WHERE id = $5`,
		user.Login,
		user.Pass,
		user.Name,
		user.PermID,
		user.ID,
	)
	if err != nil {
		return err
	}
	return nil
}

// selectUserByLogin возвращает пользователя по логину
func (db *database) selectUserByLogin(login string) (*User, error) {
	user := &User{}
	row := db.db.QueryRow(`SELECT id, login, pass,  name , id_permission FROM auth.users WHERE login = $1`, login)
	err := row.Scan(
		&user.ID,
		&user.Login,
		&user.Pass,
		&user.Name,
		&user.PermID,
	)
	if err != nil {
		return nil, err
	}
	return user, nil
}

// GetUsers возврашает всех пользователей
func (a *Auth) GetUsers() ([]*User, error) {
	return a.db.selectUsers()
}
