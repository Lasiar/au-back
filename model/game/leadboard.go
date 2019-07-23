package game

type LeaderBoard struct {
	ID       int    `json:"id"`
	Name     string `json:"name"`
	Length   int    `json:"length"`
	CountLap int    `json:"count_lap"`
}

func (db *database) selectLeaderBoard() ([]*LeaderBoard, error) {
	leaderBoards := make([]*LeaderBoard, 0)
	rows, err := db.db.Query("select id, name, count, length  from game.v_leaderboard ")
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		lb := &LeaderBoard{}
		if err := rows.Scan(&lb.ID, &lb.Name, &lb.CountLap, &lb.Length); err != nil {
			return nil, err
		}
		leaderBoards = append(leaderBoards, lb)
	}
	return leaderBoards, nil
}
