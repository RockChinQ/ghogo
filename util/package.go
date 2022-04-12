package util

import (
	"encoding/binary"
	"encoding/json"
	"net"
	"sync"
	"time"
)

type PackageIO struct {
	Connection net.Conn

	Pending         []NetPackage
	PendingMutex    *sync.Mutex
	SenderWaitGroup sync.WaitGroup
}

type NetPackage struct {
	Type      string `json:"type"`
	TimeStamp int64  `json:"time-stamp"`
	Payload   string `json:"pay-load"`
}

//Convert bytes to int
func (pio *PackageIO) ReadInt() (int, error) {
	// bytesBuffer := bytes.NewBuffer()
	var x int32
	err := binary.Read(pio.Connection, binary.BigEndian, &x)
	return int(x), err
}

func (pio *PackageIO) WriteInt(n int) error {
	x := int32(n)
	err := binary.Write(pio.Connection, binary.BigEndian, x)
	return err
}

func (pio *PackageIO) ReadJSON(obj interface{}) error {
	length, err := pio.ReadInt()
	if err != nil {
		return err
	}
	jsonBytes := make([]byte, length)
	_, err = (pio.Connection).Read(jsonBytes)

	if err != nil {
		return err
	}
	err = json.Unmarshal(jsonBytes, obj)
	return err
}

func (pio *PackageIO) WriteJSON(obj interface{}) error {
	jsonBytes, err := json.Marshal(obj)
	if err != nil {
		return err
	}
	err = pio.WriteInt(len(jsonBytes))
	if err != nil {
		return err
	}
	_, err = pio.Connection.Write(jsonBytes)
	return err
}

func (pio *PackageIO) WriteNetPackage(pkgType string, payload interface{}) error {
	bytes, err := json.Marshal(payload)
	if err != nil {
		return err
	}
	np := NetPackage{
		pkgType,
		int64(time.Now().Unix()),
		string(bytes),
	}

	return pio.WriteJSON(np)
}

func (pio *PackageIO) makeSureReady() {
	if pio.Pending == nil {
		pio.Pending = make([]NetPackage, 0)
		var mutex sync.Mutex
		pio.PendingMutex = &mutex
		go pio.sender()
	}
}

//Provides a async way to write package
func (pio *PackageIO) AppendNetPackage(pkgType string, payload interface{}) error {

	bytes, err := json.Marshal(payload)
	if err != nil {
		return err
	}
	np := NetPackage{
		pkgType,
		int64(time.Now().Unix()),
		string(bytes),
	}

	pio.Pending = append(pio.Pending, np)
	pio.FlushAsync()
	return nil
}

func (pio *PackageIO) FlushAsync() {
	pio.SenderWaitGroup.Done()
}

func (pio *PackageIO) sender() {
	for true {
		pio.PendingMutex.Lock()
		if len(pio.Pending) > 0 { //if there are,send

			nps := pio.Pending[0]
			pio.Pending = pio.Pending[1:]

			pio.PendingMutex.Unlock()

			//send
			_ = pio.WriteJSON(nps) //this error may not be process

			continue
		} else { //no pending pkg,wait
			pio.PendingMutex.Unlock()

			pio.SenderWaitGroup.Add(1)
			pio.SenderWaitGroup.Wait()
		}

	}
}

func (pio *PackageIO) Disconnect() error {
	return pio.Connection.Close()
}
