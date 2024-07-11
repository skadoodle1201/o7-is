package main

import (
	"flag"
	"fmt"
	"net"
	"os"
	"strconv"

	"github.com/codecrafters-io/redis-starter-go/internal/commands"
	"github.com/codecrafters-io/redis-starter-go/internal/tools"
)

func main() {
	fmt.Println("Logs from your program will appear here!")
	port := flag.Int("port", 6379, "The port on which the Redis server listens")
	flag.Parse()
	serve, err := net.Listen("tcp", "0.0.0.0:"+strconv.Itoa(*port))
	if err != nil {
		fmt.Println("Failed to bind to port ", *port)
		os.Exit(1)
	}

	defer serve.Close()

	fmt.Println("Listening on port ", *port)

	for {
		conn, err := serve.Accept()
		if err != nil {
			fmt.Println(err.Error())
			continue
		}

		go handleConnection(conn)
	}

}

func handleConnection(conn net.Conn) (err error) {
	defer conn.Close()
	var buf = make([]byte, 128)
	for {
		count, err := conn.Read(buf)
		if err != nil {
			return fmt.Errorf("Error reading: %v", err.Error())
		}
		// Print the received message
		readBuf := buf[:count]
		parsedData, data, err := tools.Parse(readBuf)
		// fmt.Print("Message received:", parsedData)
		if len(data) != 0 {
			return fmt.Errorf("not all data are processed, data left: %b", data)
		}
		arr, ok := parsedData.(tools.Array)
		if !ok {
			return fmt.Errorf("parsed command data should be array")
		}
		operation, ok := arr[0].(tools.BulkString)
		if !ok {
			return fmt.Errorf("operation item should be string: %+v", arr[0])
		}
		args := tools.Array{}
		if len(arr) > 1 {
			args = arr[1:]
		}
		fmt.Printf("Processing %s operation with following args %+v", operation.Value, args)

		resMessage := commands.RedisCommands(operation.Value, args)
		_, err = conn.Write([]byte(resMessage))
		if err != nil {
			return fmt.Errorf("Error writing: %v", err.Error())
		}
	}

}
