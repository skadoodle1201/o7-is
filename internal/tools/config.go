package tools

import (
	"crypto/sha1"
	"fmt"
	"net"

	"github.com/google/uuid"
)

var serverId string

type ServerConfig struct {
	port     int64
	hostName string
	id       string
	role     string
	replicas []net.Conn
}

var activeServerConf ServerConfig

func InitServerConfig(port int64, hostName string, role string) {
	id := uuid.New()
	hash := sha1.New()
	hash.Write([]byte(id.String()))
	sha1Hash := fmt.Sprintf("%x", hash.Sum(nil))
	activeServerConf = ServerConfig{
		port:     port,
		hostName: hostName,
		id:       sha1Hash,
		role:     role,
		replicas: []net.Conn{},
	}
}

func ServerUUID() string {
	return activeServerConf.id
}

func ServerPort() (int64, string) {
	return activeServerConf.port, activeServerConf.hostName
}

var masterServerConf = ServerConfig{
	port:     6379,
	hostName: "0.0.0.0",
}

func MasterPortGetter() int64 {
	return masterServerConf.port
}

func MasterHostGetter() string {
	return masterServerConf.hostName
}

func AppendNewReplicaConn(conn net.Conn) {
	activeServerConf.replicas = append(activeServerConf.replicas, conn)
}

func GetReplicaConns() (replicaConn []net.Conn) {
	return activeServerConf.replicas
}

func GetActiverServerRole() string {
	return activeServerConf.role
}
