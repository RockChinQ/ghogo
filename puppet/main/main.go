package main

import (
	"ghogo/puppet/config"
	"ghogo/puppet/netio"
	"ghogo/util"
	"sync"

	"github.com/sirupsen/logrus"
)

var wg sync.WaitGroup

func main() {
	if !util.PathExists("config.json") {
		logrus.WithFields(logrus.Fields{
			"location": "main/main.go",
		}).Warn("no config.json found,using default config")
		config.OutputPuppetConfig("config.json")
	}
	config.ResetContext()

	go netio.Connect()

	wg.Add(1)
	wg.Wait()
}
