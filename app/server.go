package main

import (
	"encoding/hex"
	"flag"
	"fmt"
	"io"
	"net"
	"os"
	"strconv"
	"strings"
	"sync"

	"github.com/codecrafters-io/redis-starter-go/internal/commands"
	serverhelpers "github.com/codecrafters-io/redis-starter-go/internal/serverHelpers"
	"github.com/codecrafters-io/redis-starter-go/internal/tools"
)

func main() {
	var wg sync.WaitGroup
	fmt.Println("Logs from your program will appear here!")
	port := flag.Int("port", int(tools.MasterPortGetter()), "The port on which the Redis server listens")
	replicaOf := flag.String("replicaof", "", "This Server is a replica of this server")
	flag.Parse()
	if *replicaOf != "" {
		hostM, portM := splitHostPort(*replicaOf)
		conn, err := net.Listen("tcp", hostM+":"+portM)
		if err != nil {
			fmt.Println("Master connection failed Already running on port ", portM)
		} else {
			conn.Close()
			portMInt, _ := strconv.Atoi(portM)
			wg.Add(1)
			go func() {
				spwanServer(portMInt, tools.MASTER_ROLE, *replicaOf)
				wg.Done()
			}()
		}

		wg.Add(1)
		go func() {
			spwanServer(*port, tools.SLAVE_ROLE, *replicaOf)
			wg.Done()
		}()

	} else {
		wg.Add(1)
		go func() {
			spwanServer(*port, tools.MASTER_ROLE, "")
			wg.Done()
		}()
	}

	wg.Wait()

}

func handleConnection(conn net.Conn, role string) {
	for {
		var buf = make([]byte, 128)
		count, err := conn.Read(buf)
		if err != nil {
			if err == io.EOF {
				return
			}
			fmt.Println("Error reading: ", err.Error())
		}
		// Print the received message
		readBuf := buf[:count]
		fmt.Println(string(readBuf))
		parsedData, data, err := tools.Parse(readBuf)
		if err != nil {
			fmt.Println("Error parsing: ", err.Error())
			continue
		}
		if len(data) != 0 {
			fmt.Println("not all data are processed, data left: ", string(data))
			continue
		}
		arr, ok := parsedData.(tools.Array)
		if !ok {
			fmt.Println("parsed command data should be array")
			continue
		}
		operation, ok := arr[0].(tools.BulkString)
		if !ok {
			fmt.Println("operation item should be string: ", arr[0])
			continue
		}
		if operation.Value == "PSYNC" {
			tools.AppendNewReplicaConn(conn)
			var emptyRDB, _ = hex.DecodeString("524544495330303131fa0972656469732d76657205372e322e30fa0a72656469732d62697473c040fa056374696d65c26d08bc65fa08757365642d6d656dc2b0c41000fa08616f662d62617365c000fff06e3bfec0ff5aa2")
			emptyRDB = append(emptyRDB, []byte("\r\n")...)
			_, rdbFileErr := conn.Write(append([]byte(fmt.Sprintf("$%d\r\n", len(emptyRDB))), emptyRDB...))
			if rdbFileErr != nil {
				fmt.Printf("Error writing: %v", rdbFileErr.Error())
				continue
			}
			continue
		}
		args := tools.Array{}
		if len(arr) > 1 {
			args = arr[1:]
		}
		fmt.Printf("Processing %s operation with following args %+v", operation.Value, args)

		resMessage := commands.RedisCommands(operation.Value, args, role)
		_, err = conn.Write([]byte(resMessage))

		if err != nil {
			fmt.Printf("Error writing: %v", err.Error())
			continue
		}
	}

}

func splitHostPort(replicaof string) (string, string) {
	parts := strings.Split(replicaof, " ")
	return parts[0], parts[1]
}

func spwanServer(port int, role string, replicaOf string) {
	serve, err := net.Listen("tcp", tools.MasterHostGetter()+":"+strconv.Itoa(port))
	if err != nil {
		fmt.Println("Failed to bind to port ", port)
		os.Exit(1)
	}
	tools.InitServerConfig(int64(port), tools.MasterHostGetter(), role)
	if role == tools.SLAVE_ROLE {
		masterConn := masterConn(replicaOf)
		serverhelpers.SendHandshakePing(masterConn)
		go handleConnection(masterConn, tools.MASTER_ROLE)
	}
	for {
		conn, err := serve.Accept()
		if err != nil {
			fmt.Println(err.Error())
			continue
		}
		go handleConnection(conn, role)
	}
}

func masterConn(replicaOf string) (masterConn net.Conn) {
	hostM, portM := splitHostPort(replicaOf)
	masterConn, err := net.Dial("tcp", hostM+":"+portM)
	if err != nil {
		fmt.Println("Master connection failed Already running on port ", portM)
		os.Exit(1)
	}
	tools.AppendNewReplicaConn(masterConn)
	return masterConn
}
