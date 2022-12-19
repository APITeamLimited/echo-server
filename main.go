package main

import (
	"fmt"

	"github.com/APITeamLimited/echo-server/server"
)

func main() {
	fmt.Printf("APITeam Ping Server Starting on port %d...\n", server.Port)

	server.Run()
}
