package network

import "ghogo/util/puppet"

type IPuppetPayload interface {
	Process(h *ConsoleHandler)
}

//Puppet login
type PuppetPayloadLogin struct { //send
	puppet.PayloadLogin
}

type PuppetPayloadSubProcess struct { //recv/send(success)
	puppet.PayloadSubProcess
}

type PuppetPayloadSubProcessOut struct { //send
	puppet.PayloadSubProcessOut
}
