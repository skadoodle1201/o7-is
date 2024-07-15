package serverhelpers

import (
	"fmt"
	"net"
	"time"
)

func SendHandshakePing(serverPort int64, serverHostName string) (err error) {
	address := fmt.Sprintf("%s:%d", serverHostName, serverPort)
	conn, err := net.Dial("tcp", address)
	if err != nil {
		return fmt.Errorf("Error HandShake With Master")
	}
	_, err = conn.Write([]byte("*1\r\n$4\r\nping\r\n"))
	time.Sleep(1 * time.Second)
	_, err = conn.Write([]byte("*3\r\n$8\r\nREPLCONF\r\n$14\r\nlistening-port\r\n$4\r\n6380\r\n"))
	time.Sleep(1 * time.Second)
	_, err = conn.Write([]byte("*3\r\n$8\r\nREPLCONF\r\n$4\r\ncapa\r\n$6\r\npsync2\r\n"))
	time.Sleep(1 * time.Second)
	_, err = conn.Write([]byte("*3\r\n$5\r\nPSYNC\r\n$1\r\n?\r\n$2\r\n-1\r\n"))

	return err
}
