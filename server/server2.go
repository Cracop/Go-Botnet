package main

import (
	"fmt"
	"net"
	"os"
	"sync"
)

var (
	bots   = make(map[net.Conn]bool) // Map to store connected clients
	botMux sync.Mutex
	// Mutex to synchronize access to the clients map
	admins    = make(map[net.Conn]bool)
	adminMux  sync.Mutex
	broadcast = make(chan string) // Channel for broadcasting messages

)

func handleConnection(conn net.Conn, broadcast chan<- string) {
	defer conn.Close()

	buffer := make([]byte, 1024)
	n, err := conn.Read(buffer)

	if err != nil {
		fmt.Println("Error reading:", err)
		return
	}

	clientType := string(buffer[:n])

	//Switch de agregarlo a que bd
	switch clientType {
	case "admin":
		adminMux.Lock()
		admins[conn] = true
		adminMux.Unlock()
		fmt.Println("admin Connected")

	case "bot":
		botMux.Lock()
		bots[conn] = true
		botMux.Unlock()
		fmt.Println("bot Connected")
	}

	for {
		n, err := conn.Read(buffer)
		receivedData := string(buffer[:n])

		if err != nil {
			//Switch para quitarlo a la bd
			switch clientType {
			case "admin":
				adminMux.Lock()
				delete(admins, conn)
				adminMux.Unlock()
				//fmt.Printf("Client disconnected: %s\n", clientAddr)
				return

			case "bot":
				botMux.Lock()
				delete(bots, conn)
				botMux.Unlock()
			}

		}
		fmt.Printf(receivedData)

	}

}

func main() {
	port := "8080"
	listener, err := net.Listen("tcp", ":"+port) //Starts listening in port

	if err != nil { //If error exits
		fmt.Println("Error:", err)
		os.Exit(1)
	}
	defer listener.Close()

	fmt.Println("Server listening on port", port)

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Error accepting connection:", err)
			continue
		}
		go handleConnection(conn, broadcast) // Handle the connection concurrently in a goroutine
	}
}
