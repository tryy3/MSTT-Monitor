package client

import (
	"fmt"
	"net"
	"os"
	"strings"

	"encoding/json"
)

var (
	version string
	updater UpdateService
)

type ErrorResponse struct {
	Error string `json:"error"`
}

// TODO: Möjligtvis printa errors till någon fil.

// StartTCPServer startar en TCP server och väntar på förfrågningar.
func StartTCPServer(connIP, connPort, connType string) {
	// Börja lyssna på en port.
	l, err := net.Listen(connType, connIP+":"+connPort)
	if err != nil {
		// Om något går fel, exit.
		fmt.Println("Error listening:", err.Error())
		os.Exit(1)
	}

	// Close när programmet är färdigt.
	defer l.Close()

	// TODO: Debug
	fmt.Println("Listening on " + connIP + ":" + connPort)
	for {
		// Vänta på inkommande requests
		conn, err := l.Accept()
		if err != nil {
			fmt.Println("Error accepting: ", err.Error())
			os.Exit(1)
		}

		// Ta hand om request i en goroutine
		go handleRequest(conn)
	}
}

// Start startar hela klienten, tar hand om olika konfigs,
// startar TCP servern osv.
func Start(v string) {
	version = v
	fmt.Printf("Starting client version %s\n", version)
	updater = UpdateService{
		Version:    version,
		Identifier: "mstt-client-windows-",
	}
	StartTCPServer("0.0.0.0", "3333", "tcp")
}

func handleRequest(conn net.Conn) {
	buf := make([]byte, 1024)
	n, err := conn.Read(buf)
	if err != nil {
		fmt.Println("Error reading:", err.Error())
		return
	}
	cmd := ParseCommand(string(buf[:n]))

	var resp interface{}

	switch strings.ToLower(cmd.Name) {
	case "check_ram":
		resp = RAMCheck(cmd)
	case "check_disc":
		resp = DiscCheck(cmd)
	case "check_cpu":
		resp = CPUCheck(cmd)
	case "uptime":
		resp = UptimeCheck(cmd)
	case "info":
		resp = InfoCheck(cmd)
	case "file":
		resp = FileCheck(cmd)
	case "update":
		resp = UpdateCheck(cmd)
	default:
		resp = ErrorResponse{Error: "Unknown command"}
	}

	respBody, err := json.Marshal(resp)
	if err != nil {
		fmt.Println("Error parsing resp:", err.Error())
		return
	}

	conn.Write(respBody)
	conn.Close()
}
