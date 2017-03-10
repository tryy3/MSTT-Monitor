package server

import (
	"errors"
	"fmt"
	"net"

	"database/sql"

	"time"

	_ "github.com/go-sql-driver/mysql"
)

var (
	db                  *sql.DB
	timestampStmt       *sql.Stmt
	insertCheckStmt     *sql.Stmt
	updatePastCheckStmt *sql.Stmt
)

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

	// Return response
	return string(reply[:n]), nil
}

// Start är huvudfunktionen för servern,
// den startar alla processer så som mysql connection,
// bygger ihop alla klienter, startar loopen för att kolla checks
func Start() {
	fmt.Println("Starting the MSTT-Monitor server")

	fmt.Println("Starting Web API server")
	go StartWebServer()

	// Starta en mysql connection
	fmt.Println("Starting mysql connection")
	var err error
	db, err = sql.Open("mysql", "root:abc123@(192.168.20.149:3306)/mstt-monitor")
	if err != nil {
		fmt.Println("Error opening database connection:", err.Error())
	}
	defer db.Close()
	fmt.Println("Mysql connection established.")

	// TODO: Bättre error handling
	// Skapa alla prepare statements
	fmt.Println("Getting the latest check information")
	timestampStmt, err = db.Prepare("SELECT timestamp FROM checks WHERE id = ?")
	if err != nil {
		panic(err.Error())
	}
	defer timestampStmt.Close()
	insertCheckStmt, err = db.Prepare("INSERT INTO checks (command_id, client_id, response, error, finished) VALUES (?, ?, ?, ?, ?)")
	if err != nil {
		panic(err.Error())
	}
	defer insertCheckStmt.Close()
	updatePastCheckStmt, err = db.Prepare("Update checks SET checked=? WHERE id = ?")
	if err != nil {
		panic(err.Error())
	}
	defer updatePastCheckStmt.Close()
	err = db.Ping()
	if err != nil {
		panic(err.Error())
	}

	// Bygg ihop klient listan
	fmt.Println("Building the clients from latest checks")
	go BuildAllClients()

	for {
		// Kolla om man fortfarande har kontakt med databasen
		err = db.Ping()
		// TODO: Försök reconnect?
		if err != nil {
			panic(err.Error())
		}

		// Kolla om det finns någon klient som måste kollas
		CheckClients()

		// Vänta en sekund innan den börjar om
		time.Sleep(1 * time.Second)
	}
}
