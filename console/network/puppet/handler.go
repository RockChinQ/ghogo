package puppet

import (
	"encoding/json"
	"ghogo/console/config/kernel"
	"ghogo/util/puppet"
	"net"
	"strconv"
	"sync"
	"time"

	"ghogo/util"

	"github.com/sirupsen/logrus"
)

var UIDIndex int32 = 0
var Handlers = make(map[int32]*PuppetHandler)
var HandlersLock sync.Mutex

type PuppetHandler struct {
	util.Handler
	UID        int32                        //UID of a puppet connection
	Name       string                       //Name for human to read
	Profile    map[string]string            //profile of this handler
	SubProcess map[string]puppet.SubProcess //subprocesses
}

/*
Fields of Profile:

Version string //version of puppet
OSName //name of puppet os
BootTimeStamp int32 //time stamp of puppet booting
InstallTimeStamp int32 //time stamp of puppet installing
Process string //uuid of puppet process
Directory string //uuid of puppet directory
Host string //uuid of puppet host
*/

const (
	STATUS_ESTABLISHED = iota
	STATUS_LOGINED
	STATUS_DISCONNECTED
)

//Wrap a new handler
func newPuppetHandler(c net.Conn) {

	pio := util.PackageIO{
		Connection: c,
	}
	handler := &PuppetHandler{
		Handler: util.Handler{
			IO:     pio,
			Status: STATUS_ESTABLISHED,
		},
		UID:        UIDIndex,
		Name:       "Unknown",
		Profile:    make(map[string]string),
		SubProcess: make(map[string]puppet.SubProcess),
	}
	HandlersLock.Lock()
	Handlers[UIDIndex] = handler
	UIDIndex++
	HandlersLock.Unlock()

	//Handshake time out
	go handler.CheckLoginTimeOut()

	go handler.Handle()
}

func (ph *PuppetHandler) CheckLoginTimeOut() {
	time.Sleep(time.Duration(kernel.GetInst().PuppetTimeOut * 1000000))
	if ph.Status == STATUS_ESTABLISHED { //有可能已经登录失败,所以忽略异常
		ph.Disconnect("login time out")
	}
}

func (ph *PuppetHandler) Handle() {
	//Read protocol passcode
	i, err := ph.IO.ReadInt()
	if err != nil {
		ph.Disconnect("read passcode:" + err.Error())
		return
	}
	if i != util.PROTOCOL_PASSCODE {
		ph.Disconnect("illegal passcode:" + strconv.Itoa(i))
		return
	}
	//passcode正确
	for {
		pack := &util.NetPackage{}
		err = ph.IO.ReadJSON(pack)
		if err != nil {
			ph.Disconnect("read net package failed:" + err.Error())
			return
		}

		//parse payload
		if ph.Status == STATUS_DISCONNECTED {
			continue
		}
		if pack.Type != "PAYLOAD_LOGIN" && ph.Status == STATUS_ESTABLISHED {
			ph.Disconnect("payload received before logined")
			return
		}

		logrus.WithFields(logrus.Fields{
			"location": "netio/puppet/handler.go",
		}).Debug("Recv: " + pack.Payload)

		//payload type
		var payload IConsolePayload
		switch pack.Type {
		case "PAYLOAD_LOGIN":
			payload = &ConsolePayloadLogin{}
		case "PAYLOAD_SUB_PROCESS":
			payload = &ConsolePayloadSubProcess{}
		case "PAYLOAD_SUB_PROCESS_OUT":
			payload = &ConsolePayloadSubProcessOut{}
		}

		err = json.Unmarshal([]byte(pack.Payload), payload)
		if err != nil {
			ph.Disconnect("parse payload failed:" + err.Error())
			return
		}
		//process payload
		payload.Process(ph)
	}
}

//Disconnect this connection,but will not delete it from handler list.
func (ph *PuppetHandler) Disconnect(reason string) {

	ph.Status = STATUS_DISCONNECTED

	ph.IO.Disconnect()

	logrus.WithFields(logrus.Fields{
		"location": "netio/puppet/handler.go",
		"puppet":   ph.Name + "," + strconv.FormatInt(int64(ph.UID), 10),
	}).Trace("Puppet disconnected:" + reason)
}
