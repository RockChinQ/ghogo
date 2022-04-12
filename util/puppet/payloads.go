package puppet

import "ghogo/util"

type IPayload interface {
	Process(h *util.Handler)
}

//Puppet login
type PayloadLogin struct { //puppet
	Name    string            `json:"name"`
	Profile map[string]string `json:"profile"`
}

type PayloadSubProcess struct { //console/puppet(success)
	SubProcess string   `josn:"sub-process"` //uuid of sub-process
	Operation  string   `json:"operation"`
	Result     string   `json:"result"`
	Content    string   `json:"content"`
	Args       []string `json:"args"`
	Decoding   string   `json:"decoding"`
}

type PayloadSubProcessOut struct { //puppet
	SubProcess string `json:"sub-process"`
	Status     string `json:"status"`
	Content    string `json:"content"`
	TimeStamp  int64  `json:"time-stamp"`
}

type PayloadSubProcessList struct { //server/client
	SubProcesses []SubProcess
}
