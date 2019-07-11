package auth

import (
	"database/sql"
)

// Permission структура для отображения данных из таблицы уровней доступа
type Permission struct {
	ID   int
	Code string
	Name sql.NullString
}

// selectPermissions возвращает список уровней доступа
func (db *database) selectPermissions() ([]*Permission, error) {
	perms := make([]*Permission, 0)
	rows, err := db.db.Query(`SELECT id,code, name FROM auth.permissions`)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		perm := &Permission{}
		if err := rows.Scan(&perm.ID, &perm.Code, &perm.Name); err != nil {
			return nil, err
		}
		perms = append(perms, perm)
	}
	return perms, nil
}

// selectPermissionByCode возвращает объект привелегии по коду
func (db *database) selectPermissionByCode(code string) (*Permission, error) {
	perm := &Permission{}
	row := db.db.QueryRow(`SELECT id,code, name FROM auth.permissions WHERE code = $1`, code)
	if err := row.Scan(&perm.ID, &perm.Code, &perm.Name); err != nil {
		return nil, err
	}
	return perm, nil
}
