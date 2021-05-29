package tunnel

import (
	"net"
	"runtime"

	"github.com/liupeidong0620/hummingbird/adapter"
)

const (
	// maxUDPQueueSize is the max number of UDP packets
	// could be buffered. if queue is full, upcoming packets
	// would be dropped util queue is ready again.
	maxUDPQueueSize = 2 << 10
)

type ProxyHandle = func(tcpConn adapter.TCPConn, updPacket adapter.UDPPacket) (net.Conn, error)

type Tunnel struct {
	numUDPWorkers int
	udpMultiQueue []chan adapter.UDPPacket
	tcpQueue      chan adapter.TCPConn

	proxyHandle ProxyHandle
}

func (nel *Tunnel) Init(proxyHandle ProxyHandle) {
	nel.numUDPWorkers = max(runtime.NumCPU(), 4 /* at least 4 workers */)

	nel.tcpQueue = make(chan adapter.TCPConn) /* unbuffered */
	nel.udpMultiQueue = make([]chan adapter.UDPPacket, 0, nel.numUDPWorkers)
	nel.proxyHandle = proxyHandle

	for i := 0; i < nel.numUDPWorkers; i++ {
		nel.udpMultiQueue = append(nel.udpMultiQueue, make(chan adapter.UDPPacket, maxUDPQueueSize))
	}

	go nel.process()
}

// Add adds tcpConn to tcpQueue.
func (nel *Tunnel) Add(conn adapter.TCPConn) {
	nel.tcpQueue <- conn
}

// AddPacket adds udpPacket to udpQueue.
func (nel *Tunnel) AddPacket(packet adapter.UDPPacket) {
	m := packet.Metadata()
	// In order to keep each packet sent in order, we
	// calculate which queue each packet should be sent
	// by src/dst info, and make sure the rest of them
	// would only be sent to the same queue.
	i := int(m.SrcPort+m.DstPort) % nel.numUDPWorkers

	select {
	case nel.udpMultiQueue[i] <- packet:
	default:
		//log.Warnf("queue is currently full, packet will be dropped")
		packet.Drop()
	}
}

func (nel *Tunnel) process() {
	for _, udpQueue := range nel.udpMultiQueue {
		queue := udpQueue
		go func() {
			for packet := range queue {
				nel.handleUDP(packet)
			}
		}()
	}

	for conn := range nel.tcpQueue {
		go nel.handleTCP(conn)
	}
}
