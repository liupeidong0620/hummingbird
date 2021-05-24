package mod

import (
	"encoding/json"
	"fmt"
	"net"

	"github.com/liupeidong0620/hummingbird/adapter"
)

type Stat int

const (
	NextStat Stat = iota
	StopStat
)

var (
	registerMods map[string]Module
)

type moduleName struct {
	Name string `json:"name"`
}

type config struct {
	Cfg []interface{} `json:"module"`
}

type ModuleCfg struct {
	cfg  config
	mods []Module
}

type Module interface {
	Config(string, int) error

	Init() error

	Name() string

	Type() string

	Index() int

	Process(adapter.TCPConn, adapter.UDPPacket) (net.Conn, Stat, error)
}

func Register(mod Module) {
	if mod != nil && mod.Name() != "" {
		if registerMods == nil {
			registerMods = make(map[string]Module, 0)
		}
		registerMods[mod.Name()] = mod
	}
}

func (cfg *ModuleCfg) Init(config []byte) error {
	if config == nil || len(config) == 0 {
		return fmt.Errorf("module config is nil.")
	}

	if len(registerMods) <= 0 {
		return fmt.Errorf("no any module.")
	}

	err := json.Unmarshal(config, &cfg.cfg)
	if err != nil {
		return err
	}

	// init module
	index := 0
	for i := 0; i < len(cfg.cfg.Cfg); i++ {
		modCfg, err := json.Marshal(cfg.cfg.Cfg[i])
		if err != nil {
			return err
		}
		modName := moduleName{}
		err = json.Unmarshal(modCfg, &modName)
		if err != nil {
			return err
		}

		if _, ok := registerMods[modName.Name]; ok {
			if cfg.mods == nil {
				cfg.mods = make([]Module, 0, len(cfg.cfg.Cfg))
			}
			// init module config
			err = registerMods[modName.Name].Config(string(modCfg), index)
			if err != nil {
				return err
			}
			index++
			cfg.mods = append(cfg.mods, registerMods[modName.Name])
		}
	}

	return nil
}

func (cfg *ModuleCfg) GetModules() []Module {
	return cfg.mods
}

func (cfg *ModuleCfg) GetModuleNum() int {
	return len(cfg.mods)
}
