package serverhelpers

import (
	"fmt"
	"net"
	"strconv"
	"time"

	"github.com/skadoodle1201/o7-is/internal/tools"
)

func SendHandshakePing(conn net.Conn) (err error) {
	pingCommand := fmt.Sprintf("*1%s%s", tools.CLRF, tools.RedisBulkString("ping").Encode())
	_, err = conn.Write([]byte(pingCommand))
	if err != nil {
		return fmt.Errorf("error handShake with master")
	}
	time.Sleep(1 * time.Second)
	activePort, _ := tools.ServerPort()
	replConfComPort := fmt.Sprintf("*3%s%s%s%s",
		tools.CLRF,
		tools.RedisBulkString("REPLCONF").Encode(),
		tools.RedisBulkString("listening-port").Encode(),
		tools.RedisBulkString(strconv.Itoa(int(activePort))).Encode(),
	)
	_, err = conn.Write([]byte(replConfComPort))
	if err != nil {
		return fmt.Errorf("error handShake with master")
	}
	time.Sleep(1 * time.Second)
	replConfComCapa := fmt.Sprintf("*3%s%s%s%s",
		tools.CLRF,
		tools.RedisBulkString("REPLCONF").Encode(),
		tools.RedisBulkString("capa").Encode(),
		tools.RedisBulkString("psync2").Encode(),
	)
	_, err = conn.Write([]byte(replConfComCapa))
	if err != nil {
		return fmt.Errorf("error handShake with master")
	}
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

func SendSetCommandToReplica(conn net.Conn, key string, val string) {
	sendSetCommand := fmt.Sprintf("*3%s%s%s%s",
		tools.CLRF,
		tools.RedisBulkString("SET").Encode(),
		tools.RedisBulkString(key).Encode(),
		tools.RedisBulkString(val).Encode(),
	)
	conn.Write([]byte(sendSetCommand))
}
