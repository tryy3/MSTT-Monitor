package server

import (
	"database/sql"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
)

func NewDatabase(proto, user, password, ip, port, database string) (db *Database, err error) {
	db = &Database{}
	db.db, err = sql.Open(proto, fmt.Sprintf("%s:%s@(%s:%s)/%s", user, password, ip, port, database))
	if err != nil {
		return
	}

	db.timestampStmt, err = db.db.Prepare("SELECT `timestamp` FROM `checks` WHERE `id`=?")
	if err != nil {
		return
	}

	db.insertCheckStmt, err = db.db.Prepare("INSERT INTO `checks` (`command_id`, `client_id`, `response`, `error`, `finished`) VALUES (?, ?, ?, ?, ?)")
	if err != nil {
		return
	}

	db.updatePastCheckStmt, err = db.db.Prepare("UPDATE `checks` SET `checked`=? WHERE `id`=?")
	if err != nil {
		return
	}

	err = db.Ping()
	return
}

type Database struct {
	db                  *sql.DB
	timestampStmt       *sql.Stmt
	insertCheckStmt     *sql.Stmt
	updatePastCheckStmt *sql.Stmt
}

func (d Database) Close() (err error) {
	err = d.CloseDatabase()
	return
}

func (d Database) transaction(transFunc Transact)

type Transact interface {
	func Prepare(query string)
	func Exec(f interface{}...)
}