package app

import (
	"flag"
	"fmt"
)

type CmdParam struct {
	MTU       int
	Device    string
	LogLevel  string
	LogFile   string
	Interface string
	Version   bool
	ModuleCfg string
	Proxy     string
	Help      bool
}

var (
	Cmd CmdParam

	moduleCfgTemp string = `{"module":[{"name":"dns"},{"name":"wss","url":["%s"]},{"name":"direct"}]}`
)

func (cmd *CmdParam) check() error {
	if cmd.ModuleCfg == "" && cmd.Proxy == "" {
		return fmt.Errorf("module and proxy param is null.")
	}
	// todo check url
	return nil
}

func init() {
	flag.IntVar(&Cmd.MTU, "mtu", 1420, "Set device maximum transmission unit (MTU)")
	flag.BoolVar(&Cmd.Version, "version", false, "Show version information and quit")
	flag.StringVar(&Cmd.Device, "device", "utun123", "Use this device.")
	flag.StringVar(&Cmd.Interface, "interface", "", "Use network INTERFACE")
	flag.StringVar(&Cmd.LogFile, "logfile", "", "Log file.")
	flag.StringVar(&Cmd.LogLevel, "loglevel", "info", "Log level [debug|info|warn|error]")
	flag.StringVar(&Cmd.ModuleCfg, "module", "", "module config file.")
	flag.StringVar(&Cmd.Proxy, "proxy", "", "proxy url.[ws://1.2.3.4:80]")
	flag.BoolVar(&Cmd.Help, "help", false, "this help.")
	flag.Parse()
}
