package auth

import (
	"database/sql"
	"fmt"
	"github/Lasiar/au-back/base"
	"log"

	_ "github.com/lib/pq" //import pg driver
)

// database структура для БД
type database struct {
	db *sql.DB
}

func (db *database) connect() error {
	dbn, err := sql.Open("postgres", fmt.Sprintf(base.GetConfig().ConnStr))
	if err == nil {
		db.db = dbn
	} else {
		log.Fatalf("can't connect database: %v", err)
		return err
	}
	if err = db.db.Ping(); err != nil {
		log.Fatalf("databse not ping: %s", err)
		return err
	}
	return nil
}
