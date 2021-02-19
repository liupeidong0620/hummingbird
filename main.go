package main

import (
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/liupeidong0620/hummingbird/app"
	"github.com/liupeidong0620/hummingbird/log"
	"github.com/liupeidong0620/hummingbird/version"
)

func main() {

	if app.Cmd.Help {
		flag.Usage()
		return
	}

	if app.Cmd.Version {
		fmt.Println(version.PrintVersion())
		return
	}
	log.Info(version.SoftName, " start ....")
	appBase := &app.App{}

	err := appBase.Init(app.Cmd)

	if err != nil {
		log.Error(err)
		return
	}

	err = appBase.Run()
	if err != nil {
		log.Error(err)
		return
	}

	defer appBase.Stop()

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
	<-sigCh

	log.Info(version.SoftName, " stop.")
}
