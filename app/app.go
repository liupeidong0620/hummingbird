package app

import (
	"net"
	"os"
	"strings"

	"github.com/liupeidong0620/hummingbird/adapter"
	"github.com/liupeidong0620/hummingbird/dialer"
	"github.com/liupeidong0620/hummingbird/link/rwc"
	"github.com/liupeidong0620/hummingbird/log"
	mod "github.com/liupeidong0620/hummingbird/module"
	"github.com/liupeidong0620/hummingbird/netstack"
	"github.com/liupeidong0620/hummingbird/tun"
	"github.com/liupeidong0620/hummingbird/tunnel"

	"gvisor.dev/gvisor/pkg/tcpip/stack"
)

type App struct {
	cmd CmdParam

	dev tun.Device

	linkEp *rwc.Endpoint
	nel    *tunnel.Tunnel
	statck *stack.Stack

	modCfg *mod.ModuleCfg
	mods   []mod.Module

	fileFp *os.File
}

func (a *App) bindToInterface(name string) {
	dialer.DialHook = dialer.DialerWithInterface(name)
	dialer.ListenPacketHook = dialer.ListenPacketWithInterface(name)
}

func (a *App) initLog() error {
	var err error
	log.Info("log init.")
	if a.cmd.LogFile != "" {
		a.fileFp, err = os.OpenFile(a.cmd.LogFile, os.O_WRONLY|os.O_CREATE, 0666)
		if err != nil {
			return err
		}

		log.SetOutput(a.fileFp)
	}

	switch strings.ToLower(a.cmd.LogLevel) {
	case "debug":
		log.SetLevel(log.LogLevelDebug)
	case "error":
		log.SetLevel(log.LogLevelError)
	case "warn":
		log.SetLevel(log.LogLevelWarn)
	default:
		log.SetLevel(log.LogLevelInfo)
	}
	log.Info("log init ok.")

	// init interface
	a.bindToInterface(a.cmd.Interface)
	return nil
}

func (a *App) Init(cmd CmdParam) error {
	// init log
	a.cmd = cmd

	a.initLog()
	// init module
	a.modCfg = &mod.ModuleCfg{}
	err := a.modCfg.Init(cmd.Module)
	if err != nil {
		return err
	}

	a.mods = a.modCfg.GetModules()

	for i := 0; i < len(a.mods); i++ {
		log.Info("load module name: ", a.mods[i].Name(), " index: ", a.mods[i].Index())
	}

	return nil
}

func (a *App) Stop() {
	if a.dev != nil {
		a.dev.Close()
	}
	if a.fileFp != nil {
		a.fileFp.Close()
	}
}

func (a *App) Run() error {
	// create tun device
	var err error
	a.dev, err = tun.OpenDevice(a.cmd.Device, a.cmd.MTU)
	if err != nil {
		return err
	}
	// create link
	a.linkEp, err = rwc.New(a.dev, uint32(a.cmd.MTU))
	if err != nil {
		return err
	}

	for i := 0; i < len(a.mods); i++ {
		a.mods[i].Init()
	}

	// create tunnel
	a.nel = &tunnel.Tunnel{}
	a.nel.Init(func(tcpConn adapter.TCPConn, updPacket adapter.UDPPacket) (net.Conn, error) {
		var targetConn net.Conn
		var err error
		var stat mod.Stat

		mods := a.mods
		for i := 0; i < len(mods); i++ {
			targetConn, stat, err = mods[i].Process(tcpConn, updPacket)
			if err != nil {
				log.Error("[", mods[i].Name(), "] process err: ", err)
			}
			if stat == mod.StopStat {
				break
			}
		}

		return targetConn, err
	})
	// netstack
	a.statck, err = netstack.NewDefaultStack(a.linkEp, a.nel.Add, a.nel.AddPacket)
	if err != nil {
		return err
	}

	return nil
}
