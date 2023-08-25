package main

import (
	"bufio"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"net"
	"os"
	"strings"
	"time"

	"github.com/joho/godotenv"
)

var (
	broadcast = make(chan string) // Channel for broadcasting messages
	conn      net.Conn
)

func hashMessage(message string) string {
	// Create a new SHA-256 hash object
	hash := sha256.New()
	// Write the input data to the hash object
	hash.Write([]byte(message))
	// Get the resulting hash as a byte slice
	hashBytes := hash.Sum(nil)
	// Convert the hash to a hexadecimal string
	hashString := hex.EncodeToString(hashBytes)

	return hashString
}

func sendMessage() {
	for {
		message := <-broadcast
		//message = hashMessage(message)
		_, err := conn.Write([]byte(message)) //Sends message to clients
		if err != nil {
			fmt.Println("Error broadcasting data:", err)
			fmt.Println("Exiting the Admin Panel.")
			close(broadcast)
			os.Exit(0)
		}
	}
}

func main() {
	// Load variables from .env file
	err := godotenv.Load()
	if err != nil {
		fmt.Println("Error loading .env file:", err)
		return
	}

	// Access environment variables
	SECRET := os.Getenv("SECRET")
	//dbPort := os.Getenv("DB_PORT")

	fmt.Println("Database Host:", SECRET)
	//fmt.Println("DatabaclientSecret := credentials[1]se Port:", dbPort)

	//Connect to the server
	serverAddr := "127.0.0.1:8080"

	go sendMessage()

	go func() {
		for {
			fmt.Print("GoBot> ")
			reader := bufio.NewReader(os.Stdin)
			command, _ := reader.ReadString('\n')
			command = strings.TrimRight(command, "\n")

			// Check if the user wants to exit
			if command == "exit" {
				fmt.Println("Exiting the Admin Panel.")
				close(broadcast)
				os.Exit(0)
			}

			broadcast <- command
		}

	}()

	for {

		if conn == nil {

			newConn, err := net.Dial("tcp", serverAddr)
			if err != nil {
				fmt.Println("Error connecting to the server:", err)
				time.Sleep(time.Second) // Wait before attempting to reconnect
				continue
			}

			fmt.Println("Connected to the server")
			conn = newConn
			defer conn.Close()

			// Send client type identifier to the server
			_, err = conn.Write([]byte("admin/"))
			if err != nil {
				fmt.Println("Error sending client type:", err)
				return
			}
			//TODO: como admin mandar un mensaje de autenticación

			_, err = conn.Write([]byte(hashMessage(SECRET)))
			if err != nil {
				fmt.Println("Error sending:", err)
				return
			}
		}

	}

}
