package network

import (
	"ghogo/puppet/config"
	"ghogo/util"
	"net"
	"strconv"
	"time"

	"github.com/sirupsen/logrus"
)

var HandlerInst *ConsoleHandler

var currentID = 0

func Connect() {
	for {
		currentID++

		conn, err := net.Dial("tcp", config.GetInst().ConsoleAddress+":"+strconv.Itoa(config.GetInst().ConsolePort))
		if err != nil { //failed to connect
			logrus.WithFields(logrus.Fields{
				"location": "netio/connector.go",
			}).Error("Failed to dial: " + config.GetInst().ConsoleAddress + ":" + strconv.Itoa(config.GetInst().ConsolePort) + "," + err.Error())
			time.Sleep(10 * time.Second)
			continue
		}

		HandlerInst = &ConsoleHandler{
			Handler: util.Handler{
				IO: util.PackageIO{
					Connection: conn,
				},
				Status: STATUS_ESTABLISHED,
			},
		}

		logrus.WithFields(logrus.Fields{
			"location": "netio/connector.go",
		}).Info("Successfully connected to: " + config.GetInst().ConsoleAddress + ":" + strconv.Itoa(config.GetInst().ConsolePort))

		//检查是否登录超时
		go HandlerInst.CheckLoginTimeOut(currentID)

		go HandlerInst.Handle()
		break
	}
}
