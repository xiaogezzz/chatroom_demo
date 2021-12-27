package message

import "encoding/json"

type Message struct {
	Uid string
	Content string
	Gravatar string
	Timestamp string
	Type string
}

func NewMessage(uid, content, gravatar, timestamp string, t string) *Message {
	return &Message{
		Uid: uid,
		Content: content,
		Gravatar: gravatar,
		Timestamp: timestamp,
		Type: t,
	}
}

func (msg *Message) Encode() ([]byte, error) {
	return json.Marshal(msg)
}