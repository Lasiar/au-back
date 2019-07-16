package game

import (
	"errors"
	"net/http"

	"github.com/Lasiar/au-back/model/game"

	"github.com/Lasiar/au-back/model/auth"
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
		user, err := web.GetUser(r)
		if err != nil {
			context.SetResponse(r, err)
			return
		}
		if err := game.GetGame().CreateSession(user.ID, resp.Length); err != nil {
			context.SetError(r, err)
			return
		}
		context.SetResponse(r, struct{}{})
	})
}

func GetSessions() http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token, err := web.GetToken(r)
		if err != nil {
			context.SetError(r, err)
			return
		}
		user, _, err := auth.GetAuth().GetSession(token)
		if err != nil {
			context.SetError(r, err)
			return
		}
		sessions, err := game.GetGame().GetSessions(user.ID)
		context.SetResponse(r, sessions)
	})
}

func Guess() http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		req := struct {
			IDSession int    `json:"id_session"`
			Guess     string `json:"guess"`
		}{}
		user, err := web.GetUser(r)
		if err != nil {
			context.SetError(r, err)
			return
		}
		err = web.ParseJSON(r, &req)
		if err != nil {
			context.SetError(r, err)
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
		respText, isValid, err := game.GetGame().Guess(req.IDSession, req.Guess)
		if err != nil {
			context.SetError(r, err)
			return
		}
		resp := struct {
			Text    string `json:"text"`
			IsValid bool   `json:"is_valid"`
		}{Text: respText, IsValid: isValid}
		context.SetResponse(r, resp)
	})
}
