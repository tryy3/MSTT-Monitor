package main

import "github.com/tryy3/MSTT-Monitor/client"

var Version string

func main() {
	client.Start(Version)
}
