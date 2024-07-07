package main

import (
	"fmt"

	// Uncomment this block to pass the first stage
	"net"
	"os"

	"github.com/codecrafters-io/redis-starter-go/internal/tools"
)

func main() {
	// You can use print statements as follows for debugging, they'll be visible when running tests.
	fmt.Println("Logs from your program will appear here!")

	// Uncomment this block to pass the first stage
	//
	serve, err := net.Listen("tcp", "0.0.0.0:6379")
	if err != nil {
		fmt.Println("Failed to bind to port 6379")
		os.Exit(1)
	}

	defer serve.Close()

	fmt.Println("Listening on port 6379")

	for {
		conn, err := serve.Accept()
		if err != nil {
			fmt.Println(err.Error())
			continue
		}

		go handleConnection(conn)
	}

}

func handleConnection(conn net.Conn) {
	defer conn.Close()
	var buf = make([]byte, 128)
	for {
		count, err := conn.Read(buf)
		if err != nil {
			fmt.Println("Error reading:", err.Error())
			break
		}
		// Print the received message
		readBuf := buf[:count]
		parsedData, data, err := tools.Parse(readBuf)
		// fmt.Print("Message received:", parsedData)
		if len(data) != 0 {
			fmt.Errorf("not all data are processed, data left: %b", data)
			break
		}
		arr, ok := parsedData.(tools.Array)
		if !ok {
			fmt.Errorf("parsed command data should be array")
			break
		}
		command, ok := arr[0].(tools.BulkString)
		if !ok {
			fmt.Errorf("command item should be string: %+v", arr[0])
			return
		}
		args := tools.Array{}
		if len(arr) > 1 {
			args = arr[1:]
		}
		fmt.Printf("Processing %s command with following args %+v", command.Value, args)

		var resMessage string
		switch command.Value {
		case "PING":
			{
				resMessage = "+PONG\r\n"
			}
		case "ECHO":
			{
				echoed, _ := args[0].(tools.BulkString)
				resMessage = fmt.Sprintf("+%v\r\n", echoed.Value)
			}
		}
		// Send a response back to the client
		_, err = conn.Write([]byte(resMessage))
		if err != nil {
			fmt.Println("Error writing:", err.Error())
			break
		}
	}

}
