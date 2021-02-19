package tunnel

import (
	"io"
	"net"
	"strconv"
	"sync"
	"time"

	"github.com/liupeidong0620/hummingbird/adapter"
	"github.com/liupeidong0620/hummingbird/common/pool"
	"github.com/liupeidong0620/hummingbird/log"
	"github.com/liupeidong0620/hummingbird/manager"
)

const (
	tcpWaitTimeout  = 5 * time.Second
	relayBufferSize = pool.RelayBufferSize
)

func (nel *Tunnel) handleTCP(localConn adapter.TCPConn) {
	defer localConn.Close()

	metadata := localConn.Metadata()
	if !metadata.Valid() {
		log.Warn("[tunnel] Metadata not valid: %#v", metadata)
		return
	}
	// module process
	targetConn, err := nel.proxyHandle(localConn, nil)
	if err != nil {
		log.Warn("[tunnel] TCP dial %s error: %v", metadata.DestinationAddress(), err)
		return
	}

	if dialerAddr, ok := targetConn.LocalAddr().(*net.TCPAddr); ok {
		metadata.MidIP = dialerAddr.IP
		metadata.MidPort = uint16(dialerAddr.Port)
	} else {
		ip, p, _ := net.SplitHostPort(targetConn.LocalAddr().String())
		port, _ := strconv.ParseUint(p, 10, 16)
		metadata.MidIP = net.ParseIP(ip)
		metadata.MidPort = uint16(port)
	}

	targetConn = manager.NewTracker(targetConn, metadata)
	defer targetConn.Close()

	log.Info("[tunnel] TCP %s <--> %s", metadata.SourceAddress(), metadata.DestinationAddress())
	relay(localConn, targetConn) /* relay connections */
}

// relay copies between left and right bidirectionally.
func relay(left, right net.Conn) {
	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		defer wg.Done()
		_ = copyBuffer(right, left) /* ignore error */
		right.SetReadDeadline(time.Now().Add(tcpWaitTimeout))
	}()

	go func() {
		defer wg.Done()
		_ = copyBuffer(left, right) /* ignore error */
		left.SetReadDeadline(time.Now().Add(tcpWaitTimeout))
	}()

	wg.Wait()
}

func copyBuffer(dst io.Writer, src io.Reader) error {
	buf := pool.Get(relayBufferSize)
	defer pool.Put(buf)

	_, err := io.CopyBuffer(dst, src, buf)
	return err
}
