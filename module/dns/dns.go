package wss

import (
	"fmt"
	"net"
	"strconv"

	"github.com/liupeidong0620/hummingbird/adapter"
	mod "github.com/liupeidong0620/hummingbird/module"
)

var (
	_defaultWss = &dns{}

	_defaultDns = map[string]bool{
		"114.114.114.114": true,
		"8.8.8.8":         true,
		":53":             true,
	}
)

func init() {
	mod.Register(_defaultWss)
}

type cfg struct {
	// module name
	Name string `json:"name"`
	// dns addr
	// default 114.114.114, 8.8.8.8
	// :53
	Url []string `json:"url"`
}

type dns struct {
	index int
	Cfg   cfg
}

func (d *dns) Config(cfg string, index int) error {
	if index < 0 {
		return fmt.Errorf("module index error.")
	}
	d.index = index
	// todo
	// parse config
	return nil
}

func (d *dns) Init() error {
	return nil
}

func (d *dns) Name() string {
	return "dns"
}

func (d *dns) Type() string {
	return "dns"
}

func (d *dns) Index() int {
	return d.index
}

func (d *dns) Process(tcpConn adapter.TCPConn, udpPacket adapter.UDPPacket) (net.Conn, mod.Stat, error) {
	var metadata *adapter.Metadata
	var ip string
	var port string

	if tcpConn != nil {
		metadata = tcpConn.Metadata()
	} else if udpPacket != nil {
		metadata = udpPacket.Metadata()

	}
	ip = metadata.DstIP.String()
	port = strconv.Itoa(int(metadata.DstPort))

	// find ip and port
	if _, ok := _defaultDns[ip]; ok {
		metadata.MidScheme = "dns"
	} else if _, ok := _defaultDns[port]; ok {
		metadata.MidScheme = "dns"
	} else if _, ok := _defaultDns[ip+":"+port]; ok {
		metadata.MidScheme = "dns"
	}

	return nil, mod.NextStat, nil
}

/*
	ToDo
	问题：
	1. dns 透传以后，会影响dns调度？用户体验会变差
	2. DoH（DNS-over-HTTPS ） 怎么处理？有tls加密，很麻烦
*/
