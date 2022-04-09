package netio

import (
	"encoding/binary"
	"encoding/json"
	"net"
)

type PackageIO struct {
	Connection net.Conn
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

func (pio *PackageIO) Disconnect() error {
	return pio.Connection.Close()
}
