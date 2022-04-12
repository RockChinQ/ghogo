package puppet

import "sync"

type SubProcess struct {
	UUID            string      `json:"uuid"`
	Buffer          string      `json:"-"` //console/puppet use buffer differently
	BufferMutex     *sync.Mutex `json:"-"` //~
	CreateTimeStamp int64       `json:"create-time-stamp"`
	Command         string      `json:"command"`
	Args            []string    `json:"args"`
}
