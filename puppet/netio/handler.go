package netio

import (
	"encoding/json"
	"ghogo/puppet/config"
	"ghogo/util"
	"time"

	"github.com/sirupsen/logrus"
)

type Handler struct {
	IO     PackageIO
	Status int
}

const (
	STATUS_ESTABLISHED = iota
	STATUS_LOGINED
	STATUS_DISCONNECTED
)

func (h *Handler) CheckLoginTimeOut(connectID int) {
	time.Sleep(6 * time.Second)
	if h.Status != STATUS_ESTABLISHED && currentID == connectID { //如果currentID!=connectID,则此协程已过期
		h.Disconnect("login time out")
	}
}

func (h *Handler) Handle() {
	err := h.IO.WriteInt(util.PROTOCOL_PASSCODE)
	if err != nil {
		h.Disconnect("failed to write passcode," + err.Error())
		go Connect()
		return
	}

	//write login package

	err = h.IO.WriteNetPackage("PAYLOAD_LOGIN", PayloadLogin{
		config.GetInst().Name,
		config.GetContextMap(),
	})
	if err != nil {
		h.Disconnect("failed to write login package," + err.Error())
		go Connect()
		return
	}

	//read loop
	for {
		pack := &NetPackage{}
		err = h.IO.ReadJSON(pack)
		if err != nil {
			h.Disconnect("read net package failed:" + err.Error())
			go Connect()
			return
		}

		//parse payload
		if h.Status == STATUS_ESTABLISHED {
			h.Disconnect("payload received before logined")
			go Connect()
			return
		}

		var payload IPayload
		switch pack.Type {
		case "PAYLOAD_SUB_PROCESS":
			payload = &PayloadSubProcess{}
		}

		err = json.Unmarshal([]byte(pack.Payload), payload)
		if err != nil {
			h.Disconnect("parse payload failed:" + err.Error())
			go Connect()
			return
		}
	}
}

func (h *Handler) Disconnect(reason string) {

	h.Status = STATUS_DISCONNECTED
	h.IO.Disconnect()

	logrus.WithFields(logrus.Fields{
		"location": "netio/handler.go",
	}).Info("Disconnected:" + reason)
}
