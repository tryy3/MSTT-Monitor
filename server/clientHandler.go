package server

import (
	"database/sql"
	"fmt"
	"strings"
	"sync"
	"time"
)

var clients *Clients // En lista av alla klienter

func CreateTimestamp(t string, next int64) (time.Time, error) {
	timestamp, err := time.ParseInLocation("2006-01-02 15:04:05", t, time.FixedZone("GMT", 0))
	if err != nil {
		return time.Time{}, err
	}
	timestamp.Add(time.Second * time.Duration(next))
	return timestamp, nil
}

// BuildAllClients byggar ihop alla klienter, ska bara köras i starten för servern.
func BuildAllClients() {
	// Skapa alla prepare statement
	getGroupStmt, err := db.Prepare("SELECT command_id, next_check, stop_error FROM groups WHERE group_name=?")
	if err != nil {
		panic(err.Error())
	}
	defer getGroupStmt.Close()
	getCheckStmt, err := db.Prepare("SELECT id, timestamp, checked, error, finished FROM checks WHERE command_id=? AND client_id=? ORDER BY timestamp DESC")
	if err != nil {
		panic(err.Error())
	}
	defer getCheckStmt.Close()
	getCommandStmt, err := db.Prepare("SELECT command FROM commands WHERE id=?")
	if err != nil {
		panic(err.Error())
	}
	defer getCommandStmt.Close()
	getClientsStmt, err := db.Prepare("SELECT id, group_names, ip FROM clients")
	if err != nil {
		panic(err.Error())
	}
	defer getClientsStmt.Close()

	// Hämta alla klienter
	clientRows, err := getClientsStmt.Query()
	if err != nil {
		panic(err.Error())
	}
	defer clientRows.Close()

	var (
		id         int
		groupNames string
		ip         string
	)

	wg := &sync.WaitGroup{}

	// Loopa igenom alla klienter
	for clientRows.Next() {
		fmt.Println("Found a client, preparing to build a client.\n")
		// Skanna igenom resultatet och skapa en enkel klient
		err := clientRows.Scan(&id, &groupNames, &ip)
		if err != nil {
			panic(err.Error())
		}

		fmt.Printf("Preparing to build client %d with IP %s.\n", id, ip)

		cl := &Client{ip: ip, clientID: id, checks: []*Check{}, groups: strings.Split(groupNames, ",")}
		clients.Add(cl)

		// Bygg ihop klienten i en goroutine
		wg.Add(1)
		go BuildClient(wg, cl, getGroupStmt, getCheckStmt, getCommandStmt)
	}
	wg.Wait()
}

// BuildClient bygger ihop en klient med all nödvändig data, commands, checks osv.
func BuildClient(wg *sync.WaitGroup, cl *Client, getGroupStmt *sql.Stmt, getCheckStmt *sql.Stmt, getCommandStmt *sql.Stmt) {
	defer wg.Done()
	fmt.Printf("Building client %d\n", cl.clientID)

	cl.Lock()
	groups := cl.groups
	cl.Unlock()

	for _, group := range groups {
		// Starta query för all grupp check info.
		groupRows, err := getGroupStmt.Query(group)
		if err != nil {
			panic(err.Error())
		}
		defer groupRows.Close()

		// Loopa igenom resultatet
		for groupRows.Next() {
			var (
				// Grupp resultat
				commandID int
				nextCheck int64
				stopError bool

				// Check resultat
				timestamp  string
				checked    bool
				checkError bool
				done       bool
				pastID     int64

				// Command resultat
				command string
			)
			err := groupRows.Scan(&commandID, &nextCheck, &stopError)
			if err != nil {
				panic(err.Error())
			}

			// Påbörja skapandet av en check.
			ch := &Check{commandID: commandID, nextCheck: nextCheck, failErr: stopError}

			// Hämta senaste check infon
			err = getCheckStmt.QueryRow(ch.commandID, cl.clientID).Scan(&pastID, &timestamp, &checked, &checkError, &done)
			if err != nil {
				if err == sql.ErrNoRows {
					// Make a "fake" check
					checked = false
					checkError = false
					done = true
					pastID = -1
					ch.nextTimestamp = time.Now()
				} else {
					panic(err.Error())
				}
			} else {
				// Konvertera timestamp till time.Time
				ch.nextTimestamp, err = CreateTimestamp(timestamp, ch.nextCheck)
				if err != nil {
					panic(err.Error())
				}
			}
			ch.checked = checked
			ch.err = checkError
			ch.pastID = pastID
			ch.done = done

			// Hämta kommandot som ska skickas
			err = getCommandStmt.QueryRow(ch.commandID).Scan(&command)
			if err != nil {
				if err == sql.ErrNoRows {
					ch.checked = true
					// Gör något här
				} else {
					panic(err.Error())
				}
			}
			ch.command = command
			cl.Add(ch)
		}
	}
	fmt.Printf("Finished building client %d.\n", cl.clientID)
}

// CheckClients kollar igenom alla klienter och ser om en check ska skickas eller inte.
func CheckClients() {
	fmt.Println("Checking clients")
	wg := &sync.WaitGroup{}
	for i := clients.Length(); i >= 0; i-- {
		cl := clients.Get(i)
		if cl == nil {
			continue
		}
		cl.Lock()
		id := cl.clientID
		cl.Unlock()
		for j := cl.Length(); j >= 0; j-- {
			ch := cl.Get(j)
			if ch == nil {
				continue
			}
			ch.Lock()
			// Om checken redan har skickats eller inte
			if ch.checked {
				ch.Unlock()
				continue
			}

			// Om nuvarande check är färdig eller inte
			if !ch.done {
				ch.Unlock()
				continue
			}

			// Om checken har haft ett error och failErr är sant
			// om det är så, skippa denna check
			if ch.err && ch.failErr {
				ch.Unlock()
				continue
			}

			// Kolla om det är dags att skicka en check eller inte
			if ch.nextTimestamp.IsZero() || time.Now().Before(ch.nextTimestamp) {
				ch.Unlock()
				continue
			}

			fmt.Printf("Starting a check for %d with client %d.\n", ch.commandID, id)
			ch.Unlock()

			// Skicka en klient check i en goroutine
			wg.Add(1)
			go SendClientCheck(wg, cl, ch)
		}
	}
	wg.Wait()
}

// SendClientCheck skickar en check till en klient
func SendClientCheck(wg *sync.WaitGroup, cl *Client, ch *Check) {
	// TODO: Bättre error handling
	// Uppdatera tidigare check så att man vet att den har kollats.
	ch.Lock()
	if ch.pastID != -1 {
		_, err := updatePastCheckStmt.Exec(true, ch.pastID)
		if err != nil {
			ch.Unlock()
			panic(err.Error())
		}
	}

	command := ch.command
	ch.Unlock()

	cl.Lock()
	ip := cl.ip
	cl.Unlock()

	// Skicka kommandot till klienten och vänta på response
	// TODO: Uppdatera denna funktion med JSON.
	resp, err := SendMessage(ip, "3333", "tcp", command)
	if err != nil || strings.Contains(resp, "Error:") {
		ch.Lock()
		ch.err = true
		ch.Unlock()
		wg.Done()
		return
	}

	ch.Lock()
	commandID := ch.commandID
	checkErr := ch.err
	ch.Unlock()

	cl.Lock()
	clientID := cl.clientID
	ch.Unlock()

	// Skapa en ny check i databasen.
	stmtResp, err := insertCheckStmt.Exec(commandID, clientID, resp, checkErr, true)
	if err != nil {
		panic(err.Error())
	}

	id, err := stmtResp.LastInsertId()
	if err != nil {
		panic(err.Error())
	}

	// Hämta timestampen från klienten och konvertera den till time.Time
	var timestamp string
	err = timestampStmt.QueryRow(id).Scan(&timestamp)
	if err != nil {
		if err == sql.ErrNoRows {
			// Gör något här
		} else {
			panic(err.Error())
		}
	}

	ch.Lock()
	nextCheck := ch.nextCheck
	ch.Unlock()

	nextTimestamp, err := CreateTimestamp(timestamp, nextCheck)
	if err != nil {
		panic(err.Error())
	}

	ch.Lock()
	ch.nextTimestamp = nextTimestamp
	ch.done = true
	ch.checked = false
	ch.pastID = id
	ch.Unlock()
	wg.Done()
}
