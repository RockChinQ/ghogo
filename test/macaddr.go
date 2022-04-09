package main

import (
	"fmt"
	"net"
)

func main() {
	MAC()
}

func MAC() string {
	interfaces, err := net.Interfaces()
	if err != nil {
		return "<UnknownMAC>"
	}

	mac := ""
	for _, inter := range interfaces {
		mac += inter.HardwareAddr.String() + ","
	}
	fmt.Println(mac)
	return mac
}
