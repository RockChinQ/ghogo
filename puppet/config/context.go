package config

import (
	"net"
	"reflect"
	"runtime"
	"time"

	"github.com/google/uuid"
)

var ctx *Context

type Context struct {
	Version          string `json:"version"`
	OSName           string `json:"os-name"`
	BootTimeStamp    int64  `json:"boot-time-stamp"`
	InstallTimeStamp int64  `json:"install-time-stamp"`
	ProcessUUID      string `json:"process"`
	DirectoryUUID    string `json:"directory"`
	HostUUID         string `json:"host"`
}

func GetContext() *Context {
	return ctx
}

func GetContextMap() map[string]string {
	return Struct2Map(*ctx)
}

func Struct2Map(obj interface{}) map[string]string {
	t := reflect.TypeOf(obj)
	v := reflect.ValueOf(obj)

	var data = make(map[string]string)
	for i := 0; i < t.NumField(); i++ {
		data[t.Field(i).Name] = v.Field(i).String()
	}
	return data
}

func ResetContext() {
	ctx = &Context{
		"alpha2022031700",
		runtime.GOOS,
		time.Now().Unix(),
		cfg.InstallTimeStamp,
		uuid.NewString(), //process
		cfg.Directory,
		MAC(),
	}
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
	return mac
}
