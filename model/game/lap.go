package game

import "time"

// Round структура для отображения роундов (попытки)
type Lap struct {
	IDSession int
	DT        time.Time
	Input     string
}

func (db *database) selectLapsBySessionID(id int) ([]*Lap, error) {
	laps := make([]*Lap, 0)
	rows, err := db.db.Query("SELECT id_session, dt, input FROM game.lap  where id_session = $1", id)
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

func (db *database) insertLap(idSession int, input string) error {
	_, err := db.db.Exec("INSERT INTO game.lap (id_session, input) VALUES ($1, $2)",
		idSession,
		input,
	)
	return err
}
