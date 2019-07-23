package game

import (
	"math/rand"
	"strings"
	"sync"
)

// Game структура для работы с игрой
type Game struct {
	db *database
}

var (
	_once    sync.Once
	_game    = new(Game)
	integers = "0123456789"
)

// GetAuth возвращает текущий объект авторизации
// вызовет панику если config и log не были установлены сетерами
func GetGame() *Game {
	_once.Do(func() {
		db := &database{}
		db.connect()
		_game = &Game{db: db}
	})
	return _game
}

func generateSecret(length int) (string, error) {
	secret := strings.Builder{}
	lenIntegers := len(integers)
	for i := 0; i < length; i++ {
		if err := secret.WriteByte(integers[rand.Intn(lenIntegers)]); err != nil {
			return "", err
		}
	}
	return secret.String(), nil
}

func (g *Game) CreateSession(idUser, length int) (*Session, error) {
	secret, err := generateSecret(length)
	if err != nil {
		return nil, err
	}
	return g.db.insertSession(idUser, secret)
}

func (g *Game) GetSessions(idUser int, completed bool) ([]*Session, error) {
	return g.db.selectSessions(idUser, completed)
}

func (g *Game) Guess(idSession int, guess string) (*Lap, bool, error) {
	session, err := g.db.selectSessionByID(idSession)
	if err != nil {
		return nil, false, err
	}
	lap, err := g.db.insertLap(idSession, guess)
	if err != nil {
		return nil, false, err
	}
	return lap, session.Secret == guess, nil
}

func (g *Game) GetLapsSorted(id int) ([]*Lap, error) {
	return g.db.selectLapsBySessionID(id)
}

func (g *Game) GetSession(id int) (*Session, error) {
	return g.db.selectSessionByID(id)
}

func (g *Game) GetLeaderBoards() ([]*LeaderBoard, error) {
	return g.db.selectLeaderBoard()
}
