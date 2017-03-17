package client

import (
	"net"
	"strings"
	"syscall"

	"encoding/json"

	"github.com/bobziuchkovski/cue"
	"github.com/bobziuchkovski/cue/collector"
)

var (
	version string
	updater UpdateService
)

var log = cue.NewLogger("client")

// ErrorResponse är en strukt för att enkelt skicka tillbaka ett error.
type ErrorResponse struct {
	Error string `json:"error"`
}

// TODO: Möjligtvis printa errors till någon fil.

// StartTCPServer startar en TCP server och väntar på förfrågningar.
func StartTCPServer(connIP, connPort, connType string) {
	tcp, err := net.Listen(connType, connIP+":"+connPort)
	if err != nil {
		log.Panic(err, "Error listening on TCP")
	}

	// Close när programmet är färdigt.
	defer tcp.Close()

	// TODO: Debug
	log.WithFields(cue.Fields{
		"IP":   connIP,
		"Port": connPort,
	}).Info("Started listening for TCP requests")
	for {
		// Vänta på inkommande requests
		conn, err := tcp.Accept()
		if err != nil {
			log.Panic(err, "Error accepting TCP request")
		}

		// Ta hand om request i en goroutine
		go handleRequest(conn)
	}
}

// Start startar hela klienten, tar hand om olika konfigs,
// startar TCP servern osv.
func Start(v string) {
	cue.CollectAsync(cue.INFO, 10000, collector.Terminal{}.New())
	cue.CollectAsync(cue.INFO, 10000, collector.File{
		Path:         "client.log",
		ReopenSignal: syscall.SIGHUP, // Om jag vill rotera logs i framtiden så kan man bara skicka en SIGHUP.
	}.New())

	version = v
	log.WithFields(cue.Fields{
		"Version": version,
	}).Info("Starting MSTT-Monitor client")

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
		log.Error(err, "Error reading TCP request")
		return
	}
	cmd := ParseCommand(string(buf[:n]))

	var resp interface{}

	switch strings.ToLower(cmd.Name) {
	case "check_memory":
		resp = MemoryCheck(cmd)
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
	case "netusage":
		resp = NetworkCheck(cmd)
	default:
		resp = ErrorResponse{Error: "Unknown command"}
	}

	respBody, err := json.Marshal(resp)
	if err != nil {
		log.Error(err, "Error parsing respBody")
		return
	}

	conn.Write(respBody)
	conn.Close()
}
