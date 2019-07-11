package auth

import (
	"time"

	"github.com/lib/pq"
)

// Session структура для отображения сессий
type Session struct {
	Token      string
	UserID     int
	LastUpdate pq.NullTime
}

func (s Session) String() string { return s.Token }

// selectSessionByToken возвращает сессию по токену
func (db *database) selectSessionByToken(token string) (*Session, error) {
	session := &Session{}
	row := db.db.QueryRow(`SELECT token, user_id, last_update FROM auth.sessions WHERE token = $1`, token)
	err := row.Scan(
		&session.Token,
		&session.UserID,
		&session.LastUpdate,
	)
	if err != nil {
		return nil, err
	}
	return session, nil
}

// selectSessionsByUserID возвращает список сессий пользователя
func (db *database) selectSessionsByUserID(id int) ([]*Session, error) {
	sessions := make([]*Session, 0)
	rows, err := db.db.Query(`SELECT token, user_id, last_update FROM auth.sessions WHERE user_id = $1`, id)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		session := &Session{}
		if err := rows.Scan(
			&session.Token,
			&session.UserID,
			&session.LastUpdate,
		); err != nil {
		} else {
			sessions = append(sessions, session)
		}
	}
	return sessions, nil
}

// insertSession добавляет новую сессию
func (db *database) insertSession(session *Session) error {
	_, err := db.db.Exec(
		`INSERT INTO auth.sessions (token, user_id, last_update) VALUES ($1, $2, $3)`,
		session.Token,
		session.UserID,
		session.LastUpdate,
	)
	if err != nil {
		return err
	}
	return nil
}

// updateSessionTime обновляет последнее время использования сессии
func (db *database) updateSessionTime(token string, lastUpdate time.Time) error {
	_, err := db.db.Exec(
		`UPDATE auth.sessions SET last_update = $1 WHERE token = $2`,
		lastUpdate,
		token,
	)
	if err != nil {
		return err
	}
	return nil
}

// deleteSessionByToken удаляет сессию по токену
func (db *database) deleteSessionByToken(token string) (int, error) {
	res, err := db.db.Exec(`DELETE FROM auth.sessions WHERE token = $1`, token)
	if err != nil {
		return 0, err
	}
	c, err := res.RowsAffected()
	if err != nil {
		return 0, err
	}
	return int(c), err

}
