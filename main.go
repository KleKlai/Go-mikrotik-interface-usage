package main

import (
	"fmt"
	"io"
	"log"
	"os"

	"github.com/joho/godotenv"
	"golang.org/x/crypto/ssh"
)

func main() {

	// Initialize env
	err := godotenv.Load(".env")

	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	// Setup SSH client configuration
	config := &ssh.ClientConfig{
		User: fmt.Sprintf(os.Getenv("USERNAME")),
		Auth: []ssh.AuthMethod{
			ssh.Password(os.Getenv("PASSWORD")),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}

	// Connect to MikroTik router
	conn, err := ssh.Dial("tcp", os.Getenv("HOST"), config)
	if err != nil {
		fmt.Println("Failed to connect:", err)
		os.Exit(1)
	}
	defer conn.Close()

	// Start a new SSH session
	session, err := conn.NewSession()
	if err != nil {
		fmt.Println("Failed to create session:", err)
		os.Exit(1)
	}
	defer session.Close()

	// Set up a pipe for the session's standard output
	stdout, err := session.StdoutPipe()
	if err != nil {
		fmt.Println("Failed to setup stdout pipe:", err)
		return
	}

	// Start the session and execute multiple commands
	err = session.Start("/interface ethernet print stats where name=Uplink-ether1-Globe-SME; /interface ethernet print stats where name=Uplink-ether2-Converge")
	if err != nil {
		fmt.Println("Failed to start session:", err)
		return
	}

	// Open the file for writing
	file, err := os.Create("result.txt")
	if err != nil {
		fmt.Println("Failed to create file:", err)
		return
	}
	defer file.Close()

	// Read the session's standard output
	io.Copy(file, stdout)
	if err != nil {
		fmt.Println("Failed to read stdout:", err)
		return
	}

	// Wait for the session to finish
	err = session.Wait()
	if err != nil {
		fmt.Println("Failed to wait:", err)
		return
	}

	// pd.Export("output", string(output))

	// // Get total rx and tx bytes for a specific interface
	// interfaceGlobe := "Uplink-ether1-Globe-SME"
	// globe := fmt.Sprintf("/interface ethernet print stats where name=%s", interfaceGlobe)

	// outputGlobe, err := session.CombinedOutput(globe)
	// if err != nil {
	// 	fmt.Println("Failed to run command:", err)
	// 	os.Exit(1)
	// }

	// interfaceConverge := "Uplink-ether2-Converge"
	// converge := fmt.Sprintf("/interface ethernet print stats where name=%s", interfaceConverge)

	// outputConverge, err := session.CombinedOutput(converge)
	// if err != nil {
	// 	fmt.Println("Failed to run command:", err)
	// 	os.Exit(1)
	// }

	// fmt.Print(string(outputGlobe))
	// fmt.Print(string(outputConverge))

	// Split output into lines
	// lines := strings.Split(string(output), "\n")

	// for _, line := range lines {

	// 	if strings.Contains(line, "rx-byte") {
	// 		// Split line into words
	// 		words := strings.Fields(line)

	// 		var combined []string

	// 		for _, num := range words[1:] {
	// 			combined = append(combined, num)
	// 		}

	// 		fmt.Println(combined)
	// 		fmt.Println(strings.Join(combined, " "))
	// 		break
	// 	}
	// }
}
