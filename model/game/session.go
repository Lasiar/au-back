package game

// Session структура для отображения сессии игр
type Session struct {
	ID     int    `json:"id"`
	IDUser int    `json:"id_user"`
	Secret string `json:"secret"`
}

func (db database) selectSessionsByUserID(idUser int) ([]*Session, error) {
	sessions := make([]*Session, 0)
	rows, err := db.db.Query("SELECT id, id_user, secret FROM game.sessions WHERE id_user = $1", idUser)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		session := &Session{}
		if err := rows.Scan(
			&session.ID,
			&session.IDUser,
			&session.Secret,
		); err != nil {
			return nil, err
		}
		sessions = append(sessions, session)
	}
	return sessions, nil
}

func (db database) selectSessionByID(idUser int) (*Session, error) {
	session := &Session{}
	err := db.db.QueryRow("SELECT id, id_user, secret FROM game.sessions WHERE id = $1", idUser).Scan(
		&session.ID,
		&session.IDUser,
		&session.Secret,
	)
	if err != nil {
		return nil, err
	}
	return session, nil
}

func (db database) insertSession(idUser int, secret string) error {
	_, err := db.db.Exec("INSERT INTO game.sessions (id_user, secret) VALUES ($1, $2)", idUser, secret)
	return err
}
