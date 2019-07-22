package game

import (
	"time"
)

// Round структура для отображения роундов (попытки)
type Lap struct {
	IDSession int       `json:"id_session"`
	DT        time.Time `json:"date_time"`
	Input     string    `json:"input"`
}

func (db *database) selectLapsBySessionID(id int) ([]*Lap, error) {
	laps := make([]*Lap, 0)
	rows, err := db.db.Query("SELECT id_session, dt, input FROM game.lap  WHERE id_session = $1 ORDER BY  dt", id)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		lap := &Lap{}
		if err := rows.Scan(&lap.IDSession,
			&lap.DT,
			&lap.Input,
		); err != nil {
			return nil, err
		}
		laps = append(laps, lap)
	}
	return laps, nil
}

func (db *database) insertLap(idSession int, input string) (*Lap, error) {
	lap := &Lap{IDSession: idSession, Input: input}
	err := db.db.QueryRow("INSERT INTO game.lap (id_session, input) VALUES ($1, $2) RETURNING dt",
		idSession,
		input,
	).Scan(&lap.DT)
	return lap, err
}
