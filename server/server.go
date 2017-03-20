package server

import (
	"errors"
	"net"
	"syscall"

	"database/sql"

	"time"

	"fmt"

	"github.com/bobziuchkovski/cue"
	"github.com/bobziuchkovski/cue/collector"
	_ "github.com/go-sql-driver/mysql"
)

var (
	db                  *sql.DB
	timestampStmt       *sql.Stmt
	insertCheckStmt     *sql.Stmt
	updatePastCheckStmt *sql.Stmt

	config *Config

	clients *Clients // En lista av alla klienter
)

var log = cue.NewLogger("server")

// SendMessage skickar ett TCP meddelande
func SendMessage(connIP, connPort, connType, message string) (string, error) {
	// Starta en kontakt med en TCP server
	conn, err := net.Dial(connType, connIP+":"+connPort)
	if err != nil {
		return "", errors.New("Failed to connect: " + err.Error())
	}
	defer conn.Close()

	// Skicka meddelandet
	if _, err = conn.Write([]byte(message)); err != nil {
		return "", errors.New("Failed to write to server: " + err.Error())
	}

	// Läs responsen
	reply := make([]byte, 1024)
	n, err := conn.Read(reply)
	if err != nil {
		return "", errors.New("Failed to write to server: " + err.Error())
	}

	return string(reply[:n]), nil
}

// Start är huvudfunktionen för servern,
// den startar alla processer så som mysql connection,
// bygger ihop alla klienter, startar loopen för att kolla checks
func Start() {
	log := cue.NewLogger("server")
	cue.CollectAsync(cue.DEBUG, 10000, collector.Terminal{}.New())
	cue.CollectAsync(cue.DEBUG, 10000, collector.File{
		Path:         "server.log",
		ReopenSignal: syscall.SIGHUP, // Om jag vill rotera logs i framtiden så kan man bara skicka en SIGHUP.
	}.New())

	config = &Config{}
	config.Load()

	clients = &Clients{}

	log.Info("Starting the MSTT-Monitor server")

	log.Info("Starting Web API server")
	go StartWebServer()

	// Starta en mysql connection
	log.Info("Starting mysql connection")
	var err error
	db, err = sql.Open(config.SQLProtocol, fmt.Sprintf("%s:%s@(%s:%s)/%s", config.SQLUser, config.SQLPassword, config.SQLIP, config.SQLPort, config.SQLDatabase))
	if err != nil {
		log.Panic(err, "Error opening database connection")
	}
	defer db.Close()
	log.Info("Mysql connection established.")

	// TODO: Bättre error handling
	// Skapa alla prepare statements
	log.Info("Skapar mysql prepare statements")

	timestampStmt, err = db.Prepare("SELECT timestamp FROM checks WHERE id = ?")
	if err != nil {
		log.Panic(err, "Error creating the timestamp prepare statement")
	}
	defer timestampStmt.Close()

	insertCheckStmt, err = db.Prepare("INSERT INTO checks (command_id, client_id, response, error, finished) VALUES (?, ?, ?, ?, ?)")
	if err != nil {
		log.Panic(err, "Error creating the insert Check prepare statement")
	}
	defer insertCheckStmt.Close()

	updatePastCheckStmt, err = db.Prepare("Update checks SET checked=? WHERE id = ?")
	if err != nil {
		log.Panic(err, "Error creating the Update Past Check prepare statement")
	}
	defer updatePastCheckStmt.Close()

	err = db.Ping()
	if err != nil {
		log.Panic(err, "Error pinging the database")
	}

	// Bygg ihop klient listan
	log.Info("Building the clients from latest checks")
	go BuildAllClients()

	for {
		// Kolla om man fortfarande har kontakt med databasen
		err = db.Ping()

		// TODO: Försök reconnect?
		if err != nil {
			log.Panic(err, "Error pinging the database")
		}

		// Kolla om det finns någon klient som måste kollas
		CheckClients()

		// Vänta en sekund innan den börjar om
		time.Sleep(time.Duration(config.Interval) * time.Second)
	}
}
