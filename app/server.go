package main

import (
	"encoding/hex"
	"flag"
	"fmt"
	"net"
	"os"
	"strconv"

	"github.com/codecrafters-io/redis-starter-go/internal/commands"
	serverhelpers "github.com/codecrafters-io/redis-starter-go/internal/serverHelpers"
	"github.com/codecrafters-io/redis-starter-go/internal/tools"
)

func main() {
	fmt.Println("Logs from your program will appear here!")
	port := flag.Int("port", int(tools.MasterPortGetter()), "The port on which the Redis server listens")
	role := flag.String("replicaof", tools.MASTER_ROLE, "The role redis server is running on")
	flag.Parse()
	serve, err := net.Listen("tcp", tools.MasterHostGetter()+":"+strconv.Itoa(*port))
	if err != nil {

		fmt.Println("Failed to bind to port ", *port)
		os.Exit(1)
	}
	roleCreated := tools.MASTER_ROLE
	if *role != tools.MASTER_ROLE {
		go serverhelpers.SendHandshakePing(tools.MasterPortGetter(), tools.MasterHostGetter())
		roleCreated = tools.SLAVE_ROLE
	}
	tools.InitServerConfig(int64(*port), tools.MasterHostGetter(), roleCreated)

	defer serve.Close()

	fmt.Println("Listening on port ", *port)

	for {
		conn, err := serve.Accept()

		if err != nil {
			fmt.Println(err.Error())
			continue
		}

		go handleConnection(conn, *role)
	}

}

func handleConnection(conn net.Conn, role string) (err error) {
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

		resMessage := commands.RedisCommands(operation.Value, args, role)
		_, err = conn.Write([]byte(resMessage))

		if operation.Value == "PSYNC" {
			tools.AppendNewReplicaConn(conn)
			var emptyRDB, _ = hex.DecodeString("524544495330303131fa0972656469732d76657205372e322e30fa0a72656469732d62697473c040fa056374696d65c26d08bc65fa08757365642d6d656dc2b0c41000fa08616f662d62617365c000fff06e3bfec0ff5aa2")
			_, err = conn.Write(append([]byte(fmt.Sprintf("$%d\r\n", len(emptyRDB))), emptyRDB...))

		}
		if err != nil {
			return fmt.Errorf("Error writing: %v", err.Error())
		}
	}

}
