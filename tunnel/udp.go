package tunnel

import (
	"errors"
	"net"
	"os"
	"strconv"
	"time"

	"github.com/liupeidong0620/hummingbird/adapter"
	"github.com/liupeidong0620/hummingbird/common/pool"
	"github.com/liupeidong0620/hummingbird/log"
	"github.com/liupeidong0620/hummingbird/manager"
	"github.com/liupeidong0620/hummingbird/nat"
)

const (
	udpTimeout    = 60 * time.Second
	udpBufferSize = (1 << 16) - 1 // largest possible UDP datagram
)

var (
	// natTable uses source udp packet information
	// as key to store destination udp packetConn.
	natTable = nat.NewTable()
)

func (nel *Tunnel) handleUDP(packet adapter.UDPPacket) {
	metadata := packet.Metadata()
	if !metadata.Valid() {
		log.Warn("[tunnel] Metadata not valid: %#v", metadata)
		return
	}
	// ToDo
	srcKey := metadata.SourceAddress()
	dstKey := metadata.DestinationAddress()

	handle := func(drop bool) bool {
		entrey := natTable.Get(srcKey)
		if entrey == nil {
			return false
		}
		pc := entrey.Get(dstKey)
		if pc != nil {
			handleUDPToRemote(packet, pc, drop)
			return true
		}
		return false
	}

	if handle(true /* drop */) {
		return
	}

	lockKey := srcKey + dstKey + "-lock"
	cond, loaded := natTable.GetOrCreateLock(lockKey)
	go func() {
		if loaded {
			cond.L.Lock()
			cond.Wait()
			handle(true) /* drop after sending data to remote */
			cond.L.Unlock()
			return
		}

		defer func() {
			natTable.Delete(lockKey)
			cond.Broadcast()
		}()

		// module process
		targetConn, err := nel.proxyHandle(nil, packet)
		if err != nil {
			defer packet.Drop()
			log.Warn("[tunnel] UDP dial %s error: %v", metadata.DestinationAddress(), err)
			return
		}

		if dialerAddr, ok := targetConn.LocalAddr().(*net.UDPAddr); ok {
			metadata.MidIP = dialerAddr.IP
			metadata.MidPort = uint16(dialerAddr.Port)
		} else if dialerAddr, ok := targetConn.LocalAddr().(*net.TCPAddr); ok {
			metadata.MidIP = dialerAddr.IP
			metadata.MidPort = uint16(dialerAddr.Port)
		} else {
			ip, p, _ := net.SplitHostPort(targetConn.LocalAddr().String())
			port, _ := strconv.ParseUint(p, 10, 16)
			metadata.MidIP = net.ParseIP(ip)
			metadata.MidPort = uint16(port)
		}

		targetConn = manager.NewTracker(targetConn, metadata)

		go func() {
			defer targetConn.Close()
			defer packet.Drop()

			defer func() {
				// clear
				entry := natTable.Get(srcKey)
				if entry == nil {
					return
				}
				entry.Delete(dstKey)
				if entry.IsEmpty() {
					natTable.Delete(srcKey)
				}
			}()

			handleUDPToLocal(packet, targetConn, udpTimeout)
		}()
		// new udp conn
		entrey := nat.NewEntry()
		entrey.Set(dstKey, targetConn)

		natTable.Set(srcKey, entrey)

		handle(false /* drop */)
	}()
}

func handleUDPToRemote(packet adapter.UDPPacket, pc net.Conn, drop bool) {
	defer func() {
		if drop {
			packet.Drop()
		}
	}()

	remote := packet.Metadata().UDPAddr()

	if _, err := pc.Write(packet.Data()); err != nil {
		log.Warn("[UDP] write to %s error: %v", remote, err)
	}

	log.Info("[UDP] %s --> %s", packet.RemoteAddr(), remote)
}

func handleUDPToLocal(packet adapter.UDPPacket, pc net.Conn, timeout time.Duration) {
	buf := pool.Get(udpBufferSize)
	defer pool.Put(buf)

	metadata := packet.Metadata()

	from := &net.UDPAddr{
		IP:   metadata.DstIP,
		Port: int(metadata.DstPort),
	}

	for /* just loop */ {
		pc.SetReadDeadline(time.Now().Add(timeout))
		n, err := pc.Read(buf)
		if err != nil {
			if !errors.Is(err, os.ErrDeadlineExceeded) /* ignore i/o timeout */ {
				log.Warn("[UDP] ReadFrom error: %v", err)
			}
			return
		}

		// write to udp socket
		if _, err := packet.WriteBack(buf[:n], from); err != nil {
			log.Warn("[UDP] write back from %s error: %v", from, err)
			return
		}

		log.Info("[UDP] %s <-- %s", packet.RemoteAddr(), from)
	}
}
