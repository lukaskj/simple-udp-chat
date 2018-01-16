package protocol

import (
	"fmt"
	// "fmt"
	"encoding/json"
	"encoding/binary"
	"../../common"
)

const (
	ACTION_DISCONNECT uint16 = iota
	ACTION_CONNECT
	ACTION_PING
	ACTION_OK
	ACTION_BROADCAST
	ACTION_WHISPER
	ACTION_MOVE
)

type Serializable interface {
	Serialize() []byte
}

type Message struct {
	Action uint16 `json:"action"`
	Body string `json:"body"`
	Token string `json:"token"`
}

func (m *Message) Serialize() ([]byte, error) {
	// totalLength := len(m.Body) + len(m.Token) + 2

	result := make([]byte, 2)
	binary.BigEndian.PutUint16(result, m.Action)

	result = append(result, []byte(m.Body)...)
	result = append(result, []byte(m.Token)...)

	

	return result, nil
}

func (m *Message) SerializeOld() ([]byte, error) {
	result, err := json.Marshal(m)
	if err != nil {
		common.CheckError(err)
		var a []byte
		return a, err
	}
	return result, nil
}

func (m *Message) ParseStr(str string) {
	
}

func (m *Message) Parse(str []byte) {
	if len(str) == 0 { return }
	fmt.Println("ACTION_CONNECT", ACTION_CONNECT)
	fmt.Println("parse", str, binary.BigEndian.Uint16(str[:2]))
	m.Action = binary.BigEndian.Uint16(str[:2])
	m.Body = string(str[2:])
	fmt.Println(m.Body)
	
}

func (m *Message) ParseOld(str []byte) {
	if len(str) == 0 { return }
	var aux map[string]string
	json.Unmarshal(str, &aux)
	// m.Action = []byte(aux["action"])[0]
	// m.Body = aux["body"]
}