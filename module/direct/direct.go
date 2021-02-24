package direct

import (
	"fmt"
	"net"
	"time"

	"github.com/liupeidong0620/hummingbird/adapter"
	"github.com/liupeidong0620/hummingbird/dialer"
	mod "github.com/liupeidong0620/hummingbird/module"
)

var (
	_defaultDirect = &direct{}

	dialTimeout time.Duration = time.Duration(10)
)

func init() {
	mod.Register(_defaultDirect)
}

type direct struct {
	index int `json:"index"`
}

func (d *direct) Config(_ string, index int) error {
	if index < 0 {
		return fmt.Errorf("module index error.")
	}
	d.index = index
	return nil
}

func (d *direct) Init() error {
	return nil
}

func (d *direct) Name() string {
	return "direct"
}

func (d *direct) Type() string {
	return "direct"
}

func (d *direct) Index() int {
	return d.index
}

func (d *direct) Process(tcpConn adapter.TCPConn, udpPacket adapter.UDPPacket) (net.Conn, mod.Stat, error) {
	var targetConn net.Conn
	var err error
	var metadata *adapter.Metadata
	var addr string

	if tcpConn != nil {
		metadata = tcpConn.Metadata()
		addr = metadata.DestinationAddress()
		targetConn, err = dialer.DialTimeout("tcp", addr, dialTimeout*time.Second)
	} else if udpPacket != nil {
		metadata = udpPacket.Metadata()
		addr = metadata.DestinationAddress()
		targetConn, err = dialer.DialTimeout("udp", addr, dialTimeout*time.Second)
	} else {
		return nil, mod.StopStat, fmt.Errorf("input param is nil")
	}

	if err != nil {
		// print error
		return nil, mod.StopStat, err
	}

	return targetConn, mod.StopStat, err
}
