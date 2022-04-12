//Handle connections from puppets
//The center controller of puppets
package puppet

import (
	"ghogo/console/config/kernel"
	"net"
	"strconv"

	"github.com/sirupsen/logrus"
)

var acceptor *PuppetAcceptor

type PuppetAcceptor struct {
	Lsn net.Listener
}

func GetAcceptorInst() *PuppetAcceptor {
	return acceptor
}

//With config in config/kernel pkg
func InitializePuppetAcceptor() error {
	lsn, error := net.Listen("tcp", kernel.GetInst().SocketAddress+":"+strconv.Itoa(kernel.GetInst().SocketPort))
	if error != nil {
		return error
	}
	acceptor = &PuppetAcceptor{
		lsn,
	}
	return nil
}

func (pa *PuppetAcceptor) Accept() {
	logrus.WithFields(logrus.Fields{
		"location": "netio/puppet/acceptor.go",
	}).Info("Start accepting.")
	for {
		c, err := pa.Lsn.Accept()
		if err != nil {
			logrus.WithFields(logrus.Fields{
				"location": "net/puppet/acceptor.go",
			}).Panic("Failed to accept connection.", err)
			break
		}
		newPuppetHandler(c)
	}
}
