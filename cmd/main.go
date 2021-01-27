package main

import (
	"fmt"
	"time"

	"github.com/liupeidong0620/hummingbird/tun"
)

func main() {
	fmt.Println("vim-go")

	dev, err := tun.OpenDevice("tun0", 1420)
	if err != nil {
		fmt.Println(err)
		return
	}
	time.Sleep(20 * time.Second)

	defer dev.Close()
}
