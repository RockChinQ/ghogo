package main

import (
	"ghogo/puppet/config"
	"ghogo/puppet/network"
	"ghogo/util"
	"io/ioutil"
	"sync"

	"github.com/sirupsen/logrus"
	log "github.com/sirupsen/logrus"
)

var wg sync.WaitGroup

func main() {
	if !util.PathExists("config.json") {
		logrus.WithFields(logrus.Fields{
			"location": "main/main.go",
		}).Warn("no config.json found,using default config")
		config.OutputPuppetConfig("config.json")
	} else {
		bytes, error := ioutil.ReadFile("config.json")
		if error != nil { //failed to read
			log.WithFields(log.Fields{
				"location": "main/main.go",
				"file":     "config.json",
			}).Panic("Failed to read file.", error)
		}

		config.LoadPuppetConfig(string(bytes))
	}

	config.ApplyGlobalConfig()

	config.ResetContext()

	go network.Connect()

	wg.Add(1)
	wg.Wait()
}
