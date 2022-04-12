package subprocess

import (
	"bufio"
	"ghogo/puppet/config"
	"ghogo/util/puppet"
	"io"
	"os/exec"
	"sync"
	"time"
)

type ProcessHandler struct {
	puppet.SubProcess

	Stdin    io.WriteCloser `json:"-"`
	Decoding string         `json:"decoding"`
	Flush    func(buf string, sb ProcessHandler)
}

var SubProcesses []ProcessHandler
var SubProcessesLock sync.Mutex

//Make a process handler inst. with provided uuid and command and args
func MakeProcessHandler(UUID string, command string, args []string, decoding string, flush func(string, ProcessHandler)) (*ProcessHandler, error) {
	var mutex sync.Mutex
	handler := ProcessHandler{
		SubProcess: puppet.SubProcess{
			UUID:            UUID,
			Buffer:          "",
			BufferMutex:     &mutex,
			CreateTimeStamp: time.Now().Unix(),
			Command:         command,
		},
		Decoding: decoding,
		Flush:    flush,
	}

	//initialize
	inst := exec.Command(command, args...)

	stdout, err := inst.StdoutPipe()
	if err != nil {
		return nil, err
	}

	stderr, err := inst.StderrPipe()
	if err != nil {
		return nil, err
	}

	//run reader
	go handler.streamReader(stdout)
	go handler.streamReader(stderr)

	stdin, err := inst.StdinPipe()
	if err != nil {
		return nil, err
	}

	handler.Stdin = stdin

	//add to list
	SubProcessesLock.Lock()

	SubProcesses = append(SubProcesses, handler)

	SubProcessesLock.Unlock()

	//start handle
	// go handler.handle()
	err = inst.Start()

	return &handler, err
}

func (sb *ProcessHandler) streamReader(s io.ReadCloser) {
	reader := bufio.NewReader(s)
	buf := make([]byte, 1024)
	for {
		length, err := reader.Read(buf)
		if err != nil {
			//TODO 处理异常
			break
		}
		str := string(buf[:length])

		//decoding

		//append
		sb.BufferMutex.Lock()
		sb.Buffer = sb.Buffer + str

		//check
		if len(sb.Buffer) > config.GetInst().SubProcessBufferSize {

			sb.Flush(sb.Buffer, *sb)

			sb.Buffer = ""
		}

		sb.BufferMutex.Unlock()
	}
}

func (sb *ProcessHandler) handle() {

}
