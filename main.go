package main

import (
	"log"
	"math/rand"
	"time"

	"github.com/znqerz/grafana-tools-cli/cmds"
)

func main() {
	rand.Seed(time.Now().UnixNano())
	cmds := cmds.NewGrafanacliCommand()
	if err := cmds.Execute(); err != nil {
		log.Fatalf("command exec failed: %v", err)
	}
}
