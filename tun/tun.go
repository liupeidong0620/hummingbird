package tun

import (
	"io"
)

const defaultScheme = "tun"
const defaultMtu = 1420

type Device interface {
	Name() string // returns the current name
	Mtu() int

	io.ReadWriteCloser
}

func OpenDevice(name string, mtu int) (Device, error) {
	return CreateTUN(name, mtu)
}

type device struct {
	name string
	mtu  int
	rwc  io.ReadWriteCloser
}

func (dev *device) Name() string {
	return dev.name
}

func (dev *device) Mtu() int {
	return dev.mtu
}

func (dev *device) Close() error {
	if dev.rwc != nil {
		return dev.rwc.Close()
	}
	return nil
}

func (dev *device) Write(packet []byte) (int, error) {
	if dev.rwc != nil {
		return dev.rwc.Write(packet)
	}
	return 0, nil
}

func (dev *device) Read(packet []byte) (n int, err error) {
	if dev.rwc != nil {
		return dev.rwc.Read(packet)
	}
	return 0, nil
}

func CreateDevice(name string, mtu int, rwc io.ReadWriteCloser) (Device, error) {
	if rwc == nil {
		return nil, nil
	}

	device := &device{
		name: name,
		mtu:  mtu,
		rwc:  rwc,
	}

	return device, nil
}
