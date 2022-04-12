package puppet

import (
	"ghogo/util/puppet"
	"sync"
	"time"

	"github.com/sirupsen/logrus"
)

func (p *ConsolePayloadLogin) Process(ph *PuppetHandler) {
	ph.Name = p.Name
	ph.Profile = p.Profile
	ph.Status = STATUS_LOGINED

	logrus.WithFields(logrus.Fields{
		"location": "netio/puppet/processors.go",
	}).Trace("Puppet logined:" + ph.Name + " from:" + ph.IO.Connection.RemoteAddr().String())

	//get subprocess list
}

func (p *ConsolePayloadSubProcess) Process(ph *PuppetHandler) {
	//operation: RUN DISC KILL
	if p.Operation == "RUN" { //create new sub process
		if p.Result == "succ" {
			sub := puppet.SubProcess{
				UUID:            p.SubProcess,
				Buffer:          "",
				BufferMutex:     &sync.Mutex{},
				CreateTimeStamp: time.Now().Unix(),
				Command:         p.Content,
				Args:            p.Args,
			}
			ph.SubProcess[p.SubProcess] = sub
		}
	}
}

func (p *ConsolePayloadSubProcessOut) Process(ph *PuppetHandler) {

}
