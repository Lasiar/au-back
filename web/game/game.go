package game

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/Lasiar/au-back/model/game"

	web "github.com/Lasiar/au-back/web/base"
	"github.com/Lasiar/au-back/web/context"
)

func CreateSession() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		resp := struct {
			Length int `json:"length"`
		}{}
		err := web.ParseJSON(r, &resp)
		if err != nil {
			context.SetError(r, err)
			return
		}
		if resp.Length == 0 {
			context.SetError(r, fmt.Errorf("%v:%v", web.ErrBadRequest, "not allow session with zero length secret"))
			return
		}
		user, err := web.GetUser(r)
		if err != nil {
			context.SetResponse(r, err)
			return
		}
		session, err := game.GetGame().CreateSession(user.ID, resp.Length)
		if err != nil {
			context.SetError(r, err)
			return
		}
		context.SetResponse(r, session)
	})
}

func GetSessions() http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		completed, err := strconv.ParseBool(r.URL.Query()["completed"][0])
		if err != nil {
			context.SetError(r, err)
			return
		}
		user, err := web.GetUser(r)
		if err != nil {
			context.SetError(r, err)
			return
		}
		sessions, err := game.GetGame().GetSessions(user.ID, completed)
		if err != nil {
			context.SetError(r, err)
			return
		}
		context.SetResponse(r, sessions)
	})
}

func Guess() http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user, err := web.GetUser(r)
		if err != nil {
			context.SetError(r, err)
			return
		}

		req := &struct {
			IDSession int    `json:"id_session"`
			Text      string `json:"text"`
		}{}
		err = web.ParseJSON(r, req)
		if err != nil {
			context.SetError(r, fmt.Errorf("%v: %v", web.ErrBadRequest, err))
			return
		}
		if req.Text == "" || req.IDSession < 1 {
			context.SetError(r, fmt.Errorf("%v: %v", web.ErrBadRequest, "text or id_session wrong"))
			return
		}
		session, err := game.GetGame().GetSession(req.IDSession)
		if err != nil {
			context.SetError(r, err)
			return
		}
		if session.IDUser != user.ID {
			// todo: union error
			context.SetError(r, errors.New("not forbidden"))
			return
		}
		lap, isValid, err := game.GetGame().Guess(req.IDSession, req.Text)
		if err != nil {
			context.SetError(r, err)
			return
		}
		resp := struct {
			*game.Lap
			IsValid bool `json:"is_valid"`
		}{Lap: lap, IsValid: isValid}
		context.SetResponse(r, resp)
	})
}

func History() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		req := struct {
			IDSession int `json:"id_session"`
		}{}
		err := web.ParseJSON(r, &req)
		if err != nil {
			context.SetError(r, err)
			return
		}
		laps, err := game.GetGame().GetLapsSorted(req.IDSession)
		if err != nil {
			context.SetError(r, err)
			return
		}

		context.SetResponse(r, laps)
	})
}

func Leaderboard() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		leaderBoards, err := game.GetGame().GetLeaderBoards()
		if err != nil {
			context.SetError(r, err)
			return
		}
		context.SetResponse(r, leaderBoards)
	})
}
