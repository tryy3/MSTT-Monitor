package main

import (
	"github.com/bobziuchkovski/cue"
	"github.com/tryy3/MSTT-Monitor/server"
)

func main() {
	s := server.Server{}
	s.Start(cue.INFO)
}
