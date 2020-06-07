package app

import "encoding/json"

type onSendMsg struct {
	Id          int    `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Email       string `json:"email"`
}

func onSendUnmarshaling(msg *[]byte) (*onSendMsg, error) {
	onSendMsgIns := new(onSendMsg)
	err := json.Unmarshal(*msg, onSendMsgIns)

	if err != nil {
		return nil, err
	}

	return onSendMsgIns, nil
}
