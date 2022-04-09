package puppet

import (
	"time"

	"github.com/sirupsen/logrus"
)

func (p *PayloadLogin) Process(ph *PuppetHandler) {
	ph.Name = p.Name
	ph.Profile = p.Profile
	ph.Status = STATUS_LOGINED

	logrus.WithFields(logrus.Fields{
		"location": "netio/puppet/processors.go",
	}).Info("Puppet logined:" + ph.Name + " from:" + ph.IO.Connection.LocalAddr().String())
}

func (p *PayloadSubProcess) Process(ph *PuppetHandler) {
	//operation: RUN DISC KILL
	if p.Operation == "RUN" { //successfully created new sub process
		sub := SubProcess{
			p.SubProcess,
			"",
			int64(time.Now().Unix()) * 1000,
			p.Content,
		}
		ph.SubProcess[p.SubProcess] = sub
	}
}

func (p *PayloadSubProcessOut) Process(ph *PuppetHandler) {

}
