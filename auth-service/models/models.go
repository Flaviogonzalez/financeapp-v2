package data

import "database/sql"

type Models struct {
	User User
}

var Db *sql.DB

func New(conn *sql.DB) Models {
	Db = conn

	return Models{
		User: User{},
	}
}
