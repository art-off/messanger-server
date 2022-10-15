package mes

const (
	TypeRegister = "Register"
	TypeText     = "Text"
	TypeMeta     = "Meta"
	TypeError    = "Error"
)

type Message struct {
	Type    string
	Payload any
}

type TextMessage struct {
	Sender string
	Text   string
}
