package main

import (
	"bufio"
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

func sendMessage() {
	for {
		message := <-broadcast
		_, err := conn.Write([]byte(message)) //Sends message to clients
		if err != nil {
			fmt.Println("Error broadcasting data:", err)
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
	dbHost := os.Getenv("DB_HOST")
	//dbPort := os.Getenv("DB_PORT")

	fmt.Println("Database Host:", dbHost)
	//fmt.Println("Database Port:", dbPort)

	//Connect to the server
	serverAddr := "127.0.0.1:8080"

	/*
		go func() {
			for {
				fmt.Print("GoBot> ")
				reader := bufio.NewReader(os.Stdin)
				command, _ := reader.ReadString('\n')
				command = strings.TrimRight(command, "\n")

				// Check if the user wants to exit
				if command == "exit" {
					fmt.Println("Exiting the client.")
					break
				}
			}
		}()
	*/
	go sendMessage()

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
		}
		fmt.Print("GoBot> ")
		reader := bufio.NewReader(os.Stdin)
		command, _ := reader.ReadString('\n')
		command = strings.TrimRight(command, "\n")

		// Check if the user wants to exit
		if command == "exit" {
			fmt.Println("Exiting the Admin Panel.")
			break
		}

		broadcast <- command

	}

	/*
		serverAddr := "192.168.1.65:8080"
		var conn net.Conn

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
			}
		}
	*/
}
