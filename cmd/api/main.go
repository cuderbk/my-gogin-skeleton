// cmd/api/main.go
package main

import (
	"log"
)

func main() {
	server, err := InitializeServer()
	if err != nil {
		log.Fatalf("failed to initialize server: %v", err)
	}

	err = server.Run()
	if err != nil {
		log.Fatalf("server run failed: %v", err)
	}
}
