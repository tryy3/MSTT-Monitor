package server

import (
	"database/sql"
	"fmt"
	"sync"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

type checkFields struct {
	ID        int64  `db:"id"`
	Timestamp string `db:"timestamp"`
	CommandID int64  `db:"command_id"`
	ClientID  int64  `db:"client_id"`
	Response  string `db:"response"`
	Checked   bool   `db:"checked"`
	Error     bool   `db:"error"`
	Finished  bool   `db:"finished"`
}

type commandFields struct {
	ID          int64  `db:"id"`
	Command     string `db:"command"`
	Name        string `db:"namn"`
	Description string `db:"description"`
	Format      string `db:"format"`
}

type groupFields struct {
	ID        int64  `db:"id"`
	CommandID int64  `db:"command_id"`
	GroupName string `db:"group_name"`
	NextCheck int    `db:"next_check"`
	StopError bool   `db:"stop_error"`
}

type clientFields struct {
	ID         int64  `db:"id"`
	GroupNames string `db:"group_names"`
	IP         string `db:"ip"`
	Name       string `db:"namn"`
}

type alertFields struct {
	ID        int64  `db:"id"`
	AlertID   int64  `db:"alert_id"`
	ClientID  int64  `db:"client_id"`
	Timestamp string `db:"timestmap"`
	Value     string `db:"value"`
}

type alertOptionFields struct {
	ID        int64  `db:"id"`
	ClientID  int64  `db:"client_id"`
	CommandID int64  `db:"command_id"`
	Alert     string `db:"alert"`
	Value     string `db:"value"`
	Count     int64  `db:"count"`
	Delay     int64  `db:"delay"`
	Service   string `db:"service"`
}

func NewDatabase(user, password, ip, port, database string) (db *Database, err error) {
	db = &Database{}
	db.db, err = sqlx.Connect("mysql", fmt.Sprintf("%s:%s@(%s:%s)/%s", user, password, ip, port, database))
	if err != nil {
		return
	}

	db.timestampStmt, err = db.db.Preparex("SELECT `timestamp` FROM `checks` WHERE `id`=? ORDER BY `timestamp` DESC")
	if err != nil {
		return
	}

	db.insertCheckStmt, err = db.db.Preparex("INSERT INTO `checks` (`command_id`, `client_id`, `response`, `error`, `finished`) VALUES (?, ?, ?, ?, ?)")
	if err != nil {
		return
	}

	db.updatePastCheckStmt, err = db.db.Preparex("UPDATE `checks` SET `checked`=? WHERE `id`=?")
	if err != nil {
		return
	}
	return
}

type Database struct {
	db                  *sqlx.DB
	timestampStmt       *sqlx.Stmt
	insertCheckStmt     *sqlx.Stmt
	updatePastCheckStmt *sqlx.Stmt
}

func (d Database) Close() (err error) {
	err = d.db.Close()
	return
}

func (d Database) GetClients() ([]*clientFields, error) {
	c := []*clientFields{}
	err := d.db.Select(&c, "SELECT * FROM `clients`")
	return c, err
}

func (d Database) GetClient(id int64) (*clientFields, error) {
	c := &clientFields{}
	stmt, err := d.db.Preparex("SELECT * FROM `clients` WHERE `id`=?")
	if err != nil {
		return c, err
	}
	err = stmt.Get(&c, id)
	return c, err
}

func (d Database) GetCommands() ([]commandFields, error) {
	c := []commandFields{}
	err := d.db.Select(&c, "SELECT * FROM `commands`")
	return c, err
}

func (d Database) GetCommand(id int64) (commandFields, error) {
	c := commandFields{}
	stmt, err := d.db.Preparex("SELECT * FROM `commands` WHERE `id`=?")
	if err != nil {
		return c, err
	}
	err = stmt.Get(&c, id)
	return c, err
}

func (d Database) GetGroups() ([]groupFields, error) {
	g := []groupFields{}
	err := d.db.Select(&g, "SELECT * FROM `groups`")
	return g, err
}

func (d Database) GetGroup(name string) ([]groupFields, error) {
	g := []groupFields{}
	stmt, err := d.db.Preparex("SELECT * FROM `groups` WHERE `group_name`=?")
	if err != nil {
		return g, err
	}
	err = stmt.Select(&g, name)
	return g, err
}

func (d Database) GetGroupByID(id int64) (groupFields, error) {
	g := groupFields{}
	stmt, err := d.db.Preparex("SELECT * FROM `groups` WHERE `id`=?")
	if err != nil {
		return g, err
	}
	err = stmt.Get(&g, id)
	return g, err
}

func (d Database) GetGroupByCommand(name string, id int64) (groupFields, error) {
	g := groupFields{}
	stmt, err := d.db.Preparex("SELECT * FROM `groups` WHERE `group_name`=? AND `command_id`=?")
	if err != nil {
		return g, err
	}
	err = stmt.Get(&g, name, id)
	return g, err
}

func (d Database) GetCheck(stmt *sqlx.Stmt, i ...interface{}) (checkFields, error) {
	c := checkFields{}
	err := stmt.Get(&c, i...)
	return c, err
}

func (d Database) GetAlertOptions() ([]alertOptionFields, error) {
	alerts := []alertOptionFields{}
	err := d.db.Select(&alerts, "SELECT * FROM `alert_options`")
	return alerts, err
}

func (d Database) GetAlertOptionsByID(id int64) (alertOptionFields, error) {
	alert := alertOptionFields{}
	stmt, err := d.db.Preparex("SELECT * FROM `alert_options` WHERE `id`=?")
	if err != nil {
		return alert, err
	}
	defer stmt.Close()
	err = stmt.Get(alert, id)
	return alert, err
}

func (d Database) GetAlert(stmt *sqlx.Stmt, i ...interface{}) (alertFields, error) {
	a := alertFields{}
	err := stmt.Get(&a, i...)
	return a, err
}

func (d Database) Prepare(query string) (*sqlx.Stmt, error) {
	return d.db.Preparex(query)
}

func (d Database) UpdatePastCheck(id int64) error {
	_, err := d.updatePastCheckStmt.Exec(true, id)
	return err
}

func (d Database) InsertCheck(i ...interface{}) (sql.Result, error) {
	return d.insertCheckStmt.Exec(i...)
}

func (d Database) InsertAlert(i ...interface{}) (alertFields, error) {
	var a alertFields
	insertStmt, err := d.db.Preparex("INSERT INTO `alerts` (`alert_id`, `client_id`, `value`) VALUES (?,?,?)")
	if err != nil {
		return a, err
	}
	defer insertStmt.Close()
	res, err := insertStmt.Exec(i...)
	if err != nil {
		return a, err
	}

	id, err := res.LastInsertId()
	if err != nil {
		return a, err
	}

	getStmt, err := d.db.Preparex("SELECT * FROM `alerts` WHERE `id`=?")
	if err != nil {
		return a, err
	}
	defer getStmt.Close()
	return d.GetAlert(getStmt, id)
}

func (d Database) GetLastCheckTime(id int64) (string, error) {
	var time string
	err := d.timestampStmt.QueryRow(id).Scan(&time)
	return time, err
}

func (d Database) GetRealGroups(groups []groupFields, cmds []Command) []*Group {
	g := []*Group{}
	for _, groupField := range groups {
		var group *Group = nil
		for _, groupList := range g {
			if groupList.GetName() == groupField.GroupName {
				group = groupList
				break
			}
		}
		if group == nil {
			group = &Group{rw: new(sync.RWMutex), Name: groupField.GroupName}
			g = append(g, group)
		}

		for _, c := range cmds {
			if c.GetID() == groupField.CommandID {
				cmd := c.Clone()
				cmd.SetGroupID(groupField.ID)
				cmd.SetNextCheck(groupField.NextCheck)
				cmd.SetStopError(groupField.StopError)
				group.AddCommand(cmd)
				break
			}
		}
	}
	return g
}

func (d Database) GetRealCommands(cmds []commandFields) []Command {
	c := []Command{}
	for _, cc := range cmds {
		c = append(c, Command{rw: new(sync.RWMutex), ID: cc.ID, Command: cc.Command})
	}
	return c
}

func (d Database) GetGroupFromName(name string) ([]*Group, error) {
	g, err := d.GetGroup(name)
	if err != nil || len(g) <= 0 {
		return nil, err
	}
	stmt, err := d.Prepare("SELECT * FROM `commands` WHERE `id`=?")
	if err != nil {
		return nil, err
	}
	c := []commandFields{}
	for _, group := range g {
		cmd := commandFields{}
		err = stmt.Get(&cmd, group.CommandID)
		if err != nil {
			return nil, err
		}
		c = append(c, cmd)
	}

	cmds := d.GetRealCommands(c)
	groups := d.GetRealGroups(g, cmds)
	return groups, err
}
