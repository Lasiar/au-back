package game

import (
	"bytes"
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

const (
	contains = "K"
	notFound = "_"
	right    = "В"
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

func (g *Game) CreateSession(idUser, length int) error {
	secret, err := generateSecret(length)
	if err != nil {
		return err
	}
	return g.db.insertSession(idUser, secret)
}

func (g *Game) GetSessions(idUser int) ([]*Session, error) {
	return g.db.selectSessionsByUserID(idUser)
}

func (g *Game) Guess(idSession int, guess string) (string, bool, error) {
	session, err := g.db.selectSessionByID(idSession)
	if err != nil {
		return "", false, err
	}
	buf := new(bytes.Buffer)
	for i, char := range guess {
		if byte(char) == session.Secret[i] {
			buf.WriteString(right)
			continue
		}
		if strings.Contains(session.Secret, string(char)) {
			buf.WriteString(contains)
		} else {
			buf.WriteString(notFound)
		}
	}
	return buf.String(), session.Secret == guess, nil
}

func (g *Game) GetSession(id int) (*Session, error) {
	return g.db.selectSessionByID(id)
}
