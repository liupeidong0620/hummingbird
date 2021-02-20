package wss

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"net"
	"net/http"
	"net/url"
	"strconv"

	"github.com/gorilla/websocket"
	"github.com/liupeidong0620/hummingbird/adapter"
	"github.com/liupeidong0620/hummingbird/dialer"
	"github.com/liupeidong0620/hummingbird/log"
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

	url []url.URL
}

func (w *wss) Config(cfg string, index int) error {
	if index < 0 {
		return fmt.Errorf("module index error.")
	}
	w.index = index

	err := json.Unmarshal([]byte(cfg), &w.Cfg)
	if err != nil {
		return err
	}

	if len(w.Cfg.Url) <= 0 {
		return fmt.Errorf("wss url is null.")
	}

	for i := 0; i < len(w.Cfg.Url); i++ {
		url, err := url.Parse(w.Cfg.Url[i])
		if err != nil {
			return err
		}
		if url.Scheme == "" || url.Host == "" {
			return fmt.Errorf("url param error.")
		}

		if url.Scheme != "ws" && url.Scheme != "wss" {
			return fmt.Errorf("url param error.")
		}
		log.Info("[wss] ws proxy url: ", url.String())
		w.url = append(w.url, *url)
	}
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

	// ToDo
	// Choose the fastest address
	randN := rand.Intn(len(w.url))

	header := http.Header{}
	//log.Info("[wss] procotol : ", metadata.Network())
	header.Add("Protocol", metadata.Network())
	// dns proxy wss server
	if metadata.MidScheme != "dns" {
		header.Add("Scheme", metadata.MidScheme)
	}
	if metadata.MidScheme == "" {
		metadata.MidScheme = w.url[randN].Scheme
		header.Add("Scheme", metadata.MidScheme)
	}

	header.Add("Destination-Address", metadata.DstIP.String())
	header.Add("Destination-Port", strconv.Itoa(int(metadata.DstPort)))
	//header.Set("User-Agent", fmt.Sprintf("%s/%s", runtime.GOOS))
	//log.Info("[wss] url len: ", len(w.url), " wss proxy: ", w.url[randN].String())
	targetConn, err = w.newWssConn(w.url[randN].String(), header)
	if err != nil {
		// print error
		return nil, mod.NextStat, err
	}

	return targetConn, mod.StopStat, nil
}

func (w *wss) newWssConn(url string, requestHeader http.Header) (net.Conn, error) {

	var wssDialer *websocket.Dialer = &websocket.Dialer{
		/*TLSClientConfig: &tls.Config{
			InsecureSkipVerify: true,
			//ServerName:         ServerName,
		},
		HandshakeTimeout: 10 * time.Second,*/
		NetDial: func(network, addr string) (net.Conn, error) {
			//log.Info("[wss] netDial: ", network, addr)
			return dialer.Dial(network, addr)
		},
	}

	wsc, resp, err := wssDialer.Dial(url, requestHeader)
	if err != nil {
		return nil, fmt.Errorf("err: %s, websocket status code: %s", err.Error(), resp.Status)
	}

	if resp.StatusCode != 101 {
		return nil, fmt.Errorf("websocket status code: %s", resp.Status)
	}

	ws := wssconn.WSConn(wsc)

	return ws, nil
}
