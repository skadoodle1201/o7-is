package tools

const (
	CLRF = "\r\n"
)

const (
	MASTER_ROLE = "master"
)

type ServerConfig struct {
	port     int64
	hostName string
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
