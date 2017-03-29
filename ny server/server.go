package server

import (
	"github.com/tryy3/MSTT-Monitor/server/models"
	"github.com/bobziuchkovski/cue"
	"log"
)

type Server struct {
	name string
	clients *models.Clients

	config *Config
	database *Database

	log *cue.Logger
}

func (s Server) Start(level cue.Level) {
	s.log := cue.NewLogger("server")
	cue.CollectAsync(level, 10000, collector.Terminal{}.New())
	cue.CollectAsync(level, 10000, collector.File{
		Path:         "server.log",
		ReopenSignal: syscall.SIGHUP, // Om jag vill rotera logs i framtiden så kan man bara skicka en SIGHUP.
	}.New())

	s.log.Info("Starting MSTT-Monitor server")

	s.log.Info("Reading config")
	config, err = NewConfig("config.json")
	if err != nil {
		s.log.Panic(err, "Something went wrong when loading config")
		return
	}

	// TODO Start Web server

	s.log.Info("Connecting to SQL database")
	s.database, err = NewDatabase(s.config.SQLProtocol, s.config.SQLUser, s.config.SQLPassword, s.config.SQLIP, s.config.SQLPort, s.config.SQLDatabase)
	if err != nil {
		s.database.Close()
		s.log.Panic(err, "Something went wrong when loading config")
		return
	}

	
}

func (s Server) Loop() {

}

func (s )