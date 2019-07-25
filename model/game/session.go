package game

// Session структура для отображения сессии игр
type Session struct {
	ID        int    `json:"id"`
	IDUser    int    `json:"id_user"`
	Secret    string `json:"secret"`
	Completed bool   `json:"completed"`
	Interval  string `json:"interval"`
	CountLap  int    `json:"count_lap"`
}

func (db database) selectSessions(idUser int, completed bool) ([]*Session, error) {
	if completed {
		//TODO: make it better, sql builder?
		return db.selectSessionCompleted(idUser)
	}
	sessions := make([]*Session, 0)
	rows, err := db.db.Query("SELECT id, id_user, secret, completed FROM game.v_sessions_completed WHERE id_user = $1 and completed = $2", idUser, completed)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		session := &Session{}
		if err := rows.Scan(
			&session.ID,
			&session.IDUser,
			&session.Secret,
			&session.Completed,
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

func (db database) insertSession(idUser int, secret string) (*Session, error) {
	session := &Session{IDUser: idUser, Secret: secret, Completed: false}
	err := db.db.QueryRow("INSERT INTO game.sessions (id_user, secret) VALUES ($1, $2) RETURNING id", idUser, secret).Scan(&session.ID)
	return session, err
}

func (db database) selectSessionCompleted(idUser int) ([]*Session, error) {
	sessions := make([]*Session, 0)
	rows, err := db.db.Query("SELECT id, id_user, secret, count, TO_CHAR(diff, 'HH24:MI:SS') as interval FROM game.v_sessions_completed_with_diff WHERE id_user = $1", idUser)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		session := &Session{}
		if err := rows.Scan(
			&session.ID,
			&session.IDUser,
			&session.Secret,
			&session.CountLap,
			&session.Interval,
		); err != nil {
			return nil, err
		}
		sessions = append(sessions, session)
	}
	return sessions, nil
}
