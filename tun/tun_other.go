// +build darwin freebsd openbsd linux windows

package tun

import (
	"golang.zx2c4.com/wireguard/tun"

	"github.com/liupeidong0620/hummingbird/common/pool"
)

//const offset = 4

type unixTun struct {
	mtu int

	device tun.Device
}

func CreateTUN(name string, n int) (Device, error) {
	device, err := tun.CreateTUN(name, n)
	if err != nil {
		return nil, err
	}

	mtu, err := device.MTU()
	if err != nil {
		return nil, err
	}

	ut := &unixTun{
		device: device,
		mtu:    mtu,
	}

	return ut, nil
}

func (t *unixTun) Read(packet []byte) (n int, err error) {
	buf := pool.Get(offset + len(packet))
	defer pool.Put(buf)

	if n, err = t.device.Read(buf, offset); err != nil {
		return
	}

	copy(packet, buf[offset:offset+n])
	return
}

func (t *unixTun) Write(packet []byte) (int, error) {
	buf := pool.Get(offset + len(packet))
	defer pool.Put(buf)

	copy(buf[offset:], packet)
	return t.device.Write(buf[:offset+len(packet)], offset)
}

func (t *unixTun) Name() string {
	name, _ := t.device.Name()
	return name
}

func (t *unixTun) Mtu() int {
	return int(t.mtu)
}

func (t *unixTun) Close() error {
	return t.device.Close()
}
