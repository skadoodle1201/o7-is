package main

import (
	"fmt"
	// Uncomment this block to pass the first stage
	"net"
	"os"
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

	conn, err := serve.Accept()
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	defer conn.Close()
	for {
		buf := make([]byte, 128)
		_, err = conn.Read(buf)
		if err != nil {
			fmt.Println(err.Error())
			os.Exit(1)
		}

		_, err = conn.Write([]byte("+PONG\r\n"))
		if err != nil {
			fmt.Println(err.Error())
			os.Exit(1)
		}
	}

}
