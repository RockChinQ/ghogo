package puppet

import "ghogo/util/puppet"

type IConsolePayload interface {
	Process(h *PuppetHandler)
}

//Puppet login
type ConsolePayloadLogin struct {
	puppet.PayloadLogin
}

type ConsolePayloadSubProcess struct {
	puppet.PayloadSubProcess
}

type ConsolePayloadSubProcessOut struct {
	puppet.PayloadSubProcessOut
}
