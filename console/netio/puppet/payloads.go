package puppet

type IPayload interface {
	Process(ph *PuppetHandler)
}

//Puppet login
type PayloadLogin struct { //recv
	Name    string            `json:"name"`
	Profile map[string]string `json:"profile"`
}

type PayloadSubProcess struct { //send/recv(succuss)
	SubProcess string `josn:"sub-process"` //uuid of sub-process
	Operation  string `json:"operation"`
	Content    string `json:"content"`
}

type PayloadSubProcessOut struct { //recv
	SubProcess string `json:"sub-process"`
	Status     string `json:"status"`
	Content    string `json:"content"`
	TimeStamp  int64  `json:"time-stamp"`
}
