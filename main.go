package main

import (
	"encoding/json"
	"github.com/gorilla/websocket"
	"net/http"
	"server/message"
	"server/register"
	"server/user"
)

func main() {
	http.HandleFunc("/", echo)
	http.ListenAndServe(":8080", nil)
}

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

var users = make(map[*websocket.Conn]user.User)

func echo(w http.ResponseWriter, r *http.Request) {
	connection, _ := upgrader.Upgrade(w, r, nil)
	defer connection.Close()

	u, err := register.RegisterUser(connection)
	if err != nil {
		return
	}
	users[connection] = *u
	defer delete(users, connection)

	for {
		mt, m, err := connection.ReadMessage()
		if err != nil || mt == websocket.CloseMessage {
			break
		}

		message := &mes.Message{}
		err = json.Unmarshal(m, message)
		if err != nil {
			break
		}

		switch message.Type {
		case mes.TypeText:
			if messageText, ok := message.Payload.(string); ok {
				go sendMessageInRoom(messageText, connection)
			} else {
				connection.WriteJSON(&mes.Message{
					Type:    mes.TypeError,
					Payload: "InvalidMessage",
				})
			}
		case mes.TypeRegister:
			connection.WriteJSON(&mes.Message{
				Type:    mes.TypeError,
				Payload: "YouAlreadyRegistered",
			})
		}
		// Если будут появляться еще типы сообщений от пользователь -- обработать
	}
}

func sendMessageInRoom(message string, fromConn *websocket.Conn) {
	senderUser, found := users[fromConn]
	if !found {
		return
	}

	for conn, u := range users {
		if u.Room != senderUser.Room {
			continue
		}

		conn.WriteJSON(&mes.Message{
			Type: mes.TypeText,
			Payload: mes.TextMessage{
				Sender: senderUser.Username,
				Text:   message,
			},
		})
	}
}
