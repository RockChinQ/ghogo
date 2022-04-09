//Main entry of ghogo console
package main

import (
	"flag"
	"ghogo/console/config/kernel"
	"ghogo/console/netio/puppet"
	"ghogo/util"
	"io/ioutil"
	"os"
	"sync"

	log "github.com/sirupsen/logrus"
)

var configFilePath string

var wg sync.WaitGroup

func main() {
	flag.StringVar(&configFilePath, "c", "kernel.json", "specific config file,load(if exist)/generate(if not exist)")

	flag.Parse()

	if configFilePath != "" {
		if util.PathExists(configFilePath) {
			log.WithFields(log.Fields{
				"location":  "main/main.go",
				"file":      configFilePath,
				"operation": "load",
			}).Info("Using config file.")
			bytes, error := ioutil.ReadFile(configFilePath)
			if error != nil { //failed to read
				log.WithFields(log.Fields{
					"location": "main/main.go",
					"file":     configFilePath,
				}).Panic("Failed to read file.", error)
			}

			error = kernel.LoadKernelConfig(string(bytes))
			if error == nil {
				log.WithFields(log.Fields{
					"location":  "main/main.go",
					"file":      configFilePath,
					"operation": "load",
				}).Info("Config file loaded.")
			} else { //failed to load
				log.WithFields(log.Fields{
					"location": "main/main.go",
					"file":     configFilePath,
				}).Panic("Failed to load file:", error)
			}
		} else { //no such file,generate
			log.WithFields(log.Fields{
				"location":  "main/main.go",
				"file":      configFilePath,
				"operation": "generate",
			}).Info("Generating config file.")
			error := kernel.OutputKernelConfig(configFilePath)
			if error == nil {
				log.WithFields(log.Fields{
					"location":  "main/main.go",
					"file":      configFilePath,
					"operation": "generate",
				}).Info("Config file generated.")
				os.Exit(0)
			} else { //failed to generate
				log.WithFields(log.Fields{
					"location": "main/main.go",
					"file":     configFilePath,
				}).Panic("Failed to generate file:", error)
			}
		}
	}

	//parse params provided by command.
	flag.StringVar(&kernel.GetInst().SocketAddress, "sa", kernel.GetInst().SocketAddress, "socket listener address")
	flag.IntVar(&kernel.GetInst().SocketPort, "sp", kernel.GetInst().SocketPort, "socket listener port")
	flag.IntVar(&kernel.GetInst().RuntimeMode, "mode", kernel.GetInst().RuntimeMode, "runtime mode:0(Debug) 1(Normal)")
	flag.Parse()

	//Make service

	error := puppet.InitializePuppetAcceptor()
	if error != nil {
		log.WithFields(log.Fields{
			"location": "main/main.go",
		}).Error("Failed to initialize puppet acceptor.", error)
	} else {
		go puppet.GetAcceptorInst().Accept()
	}

	wg.Add(1)
	wg.Wait()
}
