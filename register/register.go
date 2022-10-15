package register

import (
	"errors"
	"github.com/gorilla/websocket"
	"server/message"
	"server/user"
)

var nonRegisterFirstMessageError = errors.New("NonRegisterFirstMessage")

func RegisterUser(conn *websocket.Conn) (*user.User, error) {
	u, err := readUserFromRegisterMessage(conn)

	for err != nil {
		conn.WriteJSON(&mes.Message{
			Type:    mes.TypeError,
			Payload: err.Error(),
		})
		u, err = readUserFromRegisterMessage(conn)
	}

	conn.WriteJSON(&mes.Message{
		Type:    mes.TypeMeta,
		Payload: "Registered",
	})

	return u, nil
}

func readUserFromRegisterMessage(conn *websocket.Conn) (*user.User, error) {
	message := &mes.Message{}
	err := conn.ReadJSON(message)

	if err != nil {
		return nil, err
	}

	if message.Type != mes.TypeRegister {
		return nil, nonRegisterFirstMessageError
	}

	if userInfoMap, ok := message.Payload.(map[string]interface{}); ok {
		return user.UserFromMap(userInfoMap)
	} else {
		return nil, errors.New("InvalidPayload")
	}
}
