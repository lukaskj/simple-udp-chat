package protocol

import (
   "errors"
   "encoding/json"
   "encoding/binary"
   "../../common"
)

const (
   ACTION_DISCONNECT uint16 = iota
   ACTION_CONNECT
   ACTION_PING
   ACTION_OK
   ACTION_ERROR
   ACTION_BROADCAST
   ACTION_WHISPER
   ACTION_MOVE
)

type Serializable interface {
   Serialize() []byte
}

type Message struct {
   Action uint16
   Body string
   Token string
   Payload string
}

func (m *Message) Serialize() ([]byte) {
   // totalLength := len(m.Body) + len(m.Token) + 2

   result := make([]byte, 2)
   binary.BigEndian.PutUint16(result, m.Action)

   token := make([]byte, 36)
   if len(m.Token) != 36 {
      m.Token = string(token)
   }

   result = append(result, []byte(m.Token)...)
   result = append(result, []byte(m.Body)...)

   return result
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

func (m *Message) Parse(str []byte) error {
   if len(str) == 0 { return errors.New("No message to parse") }
   m.Action = binary.BigEndian.Uint16(str[:2])
   m.Token = string(str[2:38])
   m.Body = string(str[38:])
   return nil
}


func (m *Message) ParseOld(str []byte) {
   if len(str) == 0 { return }
   var aux map[string]string
   json.Unmarshal(str, &aux)
   // m.Action = []byte(aux["action"])[0]
   // m.Body = aux["body"]
}