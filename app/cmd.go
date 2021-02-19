package app

import "flag"

type CmdParam struct {
	MTU       int
	Device    string
	LogLevel  string
	LogFile   string
	Interface string
	Version   bool
	Module    string
	Help      bool
}

var (
	Cmd CmdParam

	moduleCfg string = `{"module":[{"name":"dns"},{"name":"wss"},{"name":"direct"}]}`
)

func init() {
	flag.IntVar(&Cmd.MTU, "mtu", 1420, "Set device maximum transmission unit (MTU)")
	flag.BoolVar(&Cmd.Version, "version", false, "Show version information and quit")
	flag.StringVar(&Cmd.Device, "device", "utun123", "Use this device.")
	flag.StringVar(&Cmd.Interface, "interface", "", "Use network INTERFACE")
	flag.StringVar(&Cmd.LogFile, "logfile", "", "Log file.")
	flag.StringVar(&Cmd.LogLevel, "loglevel", "info", "Log level [debug|info|warn|error]")
	flag.StringVar(&Cmd.Module, "module", moduleCfg, "module config.")
	flag.BoolVar(&Cmd.Help, "help", false, "this help.")
	flag.Parse()
}
