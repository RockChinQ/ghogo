package netio

type IPayload interface {
	Process(h *Handler)
}

//Puppet login
type PayloadLogin struct { //send
	Name    string            `json:"name"`
	Profile map[string]string `json:"profile"`
}

type PayloadSubProcess struct { //recv/send(success)
	SubProcess string `josn:"sub-process"` //uuid of sub-process
	Operation  string `json:"operation"`
	Content    string `json:"content"`
}

type PayloadSubProcessOut struct { //send
	SubProcess string `json:"sub-process"`
	Status     string `json:"status"`
	Content    string `json:"content"`
	TimeStamp  int64  `json:"time-stamp"`
}
