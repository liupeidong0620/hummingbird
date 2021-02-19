package wss

import (
	"crypto/tls"
	"fmt"
	"net"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/websocket"
	"github.com/liupeidong0620/hummingbird/adapter"
	"github.com/liupeidong0620/hummingbird/dialer"
	mod "github.com/liupeidong0620/hummingbird/module"
	"github.com/liupeidong0620/hummingbird/module/wss/wssconn"
)

const (
	tlsTimeout int = 10
)

var (
	_defaultWss = &wss{}
)

func init() {
	mod.Register(_defaultWss)
}

type cfg struct {
	// module name
	Name string `json:"name"`
	// Proxy url
	Url []string `json:"url"`
}

type wss struct {
	index int
	Cfg   cfg
}

func (w *wss) Config(cfg string, index int) error {
	if index < 0 {
		return fmt.Errorf("module index error.")
	}
	w.index = index
	return nil
}

func (w *wss) Init() error {
	return nil
}

func (w *wss) Name() string {
	return "wss"
}

func (w *wss) Type() string {
	return "wss"
}

func (w *wss) Index() int {
	return w.index
}

func (w *wss) Process(tcpConn adapter.TCPConn, udpPacket adapter.UDPPacket) (net.Conn, mod.Stat, error) {
	var targetConn net.Conn
	var err error
	var metadata *adapter.Metadata

	if tcpConn != nil {
		metadata = tcpConn.Metadata()
	} else if udpPacket != nil {
		metadata = udpPacket.Metadata()
	} else {
		return nil, mod.NextStat, fmt.Errorf("input param is nil")
	}

	header := http.Header{}
	header.Add("Protocol", metadata.Network())
	// dns proxy wss server
	if metadata.MidScheme != "" {
		header.Add("Scheme", metadata.MidScheme)
	}

	header.Add("Destination-Address", metadata.DstIP.String())
	header.Add("Destination-Port", strconv.Itoa(int(metadata.DstPort)))
	//header.Set("User-Agent", fmt.Sprintf("%s/%s", runtime.GOOS))

	// ToDo
	// Choose the fastest address
	targetConn, err = w.newWssConn(w.Cfg.Url[0], header)
	if err != nil {
		// print error
		return nil, mod.NextStat, err
	}

	return targetConn, mod.NextStat, nil
}

func (w *wss) newWssConn(addr string, requestHeader http.Header) (net.Conn, error) {

	var wssDialer *websocket.Dialer = &websocket.Dialer{
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: true,
			//ServerName:         ServerName,
		},
		HandshakeTimeout: 10 * time.Second,
		NetDial: func(network, addr string) (net.Conn, error) {
			return dialer.Dial(network, addr)
		},
	}

	wsc, resp, err := wssDialer.Dial(addr, nil)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != 101 {
		return nil, err
	}

	ws := wssconn.WSConn(wsc)

	return ws, nil
}
