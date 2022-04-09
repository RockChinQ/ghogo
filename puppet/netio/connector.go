package netio

import (
	"ghogo/puppet/config"
	"net"
	"strconv"
	"time"

	"github.com/sirupsen/logrus"
)

var handler *Handler

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

		handler = &Handler{
			PackageIO{
				Connection: conn,
			},
			STATUS_ESTABLISHED,
		}

		logrus.WithFields(logrus.Fields{
			"location": "netio/connector.go",
		}).Info("Successfully connected to: " + config.GetInst().ConsoleAddress + ":" + strconv.Itoa(config.GetInst().ConsolePort))

		go handler.CheckLoginTimeOut(currentID)
		go handler.Handle()
		break
	}
}
