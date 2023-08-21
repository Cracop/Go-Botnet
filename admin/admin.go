package main

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

func main() {
	// Load variables from .env file
	err := godotenv.Load()
	if err != nil {
		fmt.Println("Error loading .env file:", err)
		return
	}

	// Access environment variables
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")

	fmt.Println("Database Host:", dbHost)
	fmt.Println("Database Port:", dbPort)
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
