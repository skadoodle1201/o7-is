package serverhelpers

import (
	"fmt"
	"net"
	"strconv"
	"time"

	"github.com/codecrafters-io/redis-starter-go/internal/tools"
)

func SendHandshakePing(serverPort int64, serverHostName string) (err error) {
	address := fmt.Sprintf("%s:%d", serverHostName, serverPort)
	conn, err := net.Dial("tcp", address)
	if err != nil {
		return fmt.Errorf("Error HandShake With Master")
	}

	pingCommand := fmt.Sprintf("*1%s%s", tools.CLRF, tools.RedisBulkString("ping").Encode())
	_, err = conn.Write([]byte(pingCommand))
	time.Sleep(1 * time.Second)
	activePort, _ := tools.ServerPort()
	replConfComPort := fmt.Sprintf("*3%s%s%s%s",
		tools.CLRF,
		tools.RedisBulkString("REPLCONF").Encode(),
		tools.RedisBulkString("listening-port").Encode(),
		tools.RedisBulkString(strconv.Itoa(int(activePort))).Encode(),
	)
	_, err = conn.Write([]byte(replConfComPort))
	time.Sleep(1 * time.Second)
	replConfComCapa := fmt.Sprintf("*3%s%s%s%s",
		tools.CLRF,
		tools.RedisBulkString("REPLCONF").Encode(),
		tools.RedisBulkString("capa").Encode(),
		tools.RedisBulkString("psync2").Encode(),
	)
	_, err = conn.Write([]byte(replConfComCapa))
	time.Sleep(1 * time.Second)
	psyncConfCom := fmt.Sprintf("*3%s%s%s%s",
		tools.CLRF,
		tools.RedisBulkString("PSYNC").Encode(),
		tools.RedisBulkString("?").Encode(),
		tools.RedisBulkString("-1").Encode(),
	)
	_, err = conn.Write([]byte(psyncConfCom))

	return err
}
