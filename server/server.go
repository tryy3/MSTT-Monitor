package server

import (
	"syscall"

	"time"

	"github.com/bobziuchkovski/cue"
	"github.com/bobziuchkovski/cue/collector"
)

type Server struct {
	handler *Handler

	config   *Config
	database *Database

	log cue.Logger
}

func (s *Server) Start(level cue.Level) {
	s.log = cue.NewLogger("server")
	cue.CollectAsync(level, 10000, collector.Terminal{}.New())
	cue.CollectAsync(level, 10000, collector.File{
		Path:         "server.log",
		ReopenSignal: syscall.SIGHUP, // Om jag vill rotera logs i framtiden så kan man bara skicka en SIGHUP.
	}.New())

	s.log.Info("Starting MSTT-Monitor server")

	s.log.Info("Reading config")
	config, err := NewConfig("config.json")
	if err != nil {
		s.log.Panic(err, "Something went wrong when loading config")
		return
	}
	s.config = config

	// TODO Start Web server

	s.log.Info("Connecting to SQL database")
	s.database, err = NewDatabase(s.config.SQLUser, s.config.SQLPassword, s.config.SQLIP, s.config.SQLPort, s.config.SQLDatabase)
	if err != nil {
		s.database.Close()
		s.log.Panic(err, "Something went wrong when loading config")
		return
	}

	s.log.Info("Creating all clients")
	handler, err := NewHandler(s.database)
	if err != nil {
		s.database.Close()
		s.log.Panic(err, "Something went wrong when creating all the clients")
		return
	}
	s.handler = handler
	s.log.Info("Finished creating all clients")

	s.log.Info("Starting web API")
	API := HTTPServer{Server: s, Handlers: map[string]APIHandler{}}
	go API.Start()
	s.log.Info("Web API started")

	s.log.Info("Starting loop")
	for {
		go s.Loop()
		time.Sleep(time.Second * time.Duration(s.config.Interval))
	}
}

func (s *Server) Loop() {
	for cl := range s.handler.IterClients() {
		go cl.Check(s)
	}
}

func (s Server) GetLogger() cue.Logger {
	return s.log
}

func (s Server) GetDatabase() *Database {
	return s.database
}

func (s Server) GetConfig() *Config {
	return s.config
}

func (s Server) GetHandler() *Handler {
	return s.handler
}
