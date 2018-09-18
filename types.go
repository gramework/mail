package mail

type Message interface {
	// Meta returns message meta information
	Meta() *MessageMeta
	// Fields returns basic message information
	Fields() *BasicFields
	// Source returns prepared message source
	Source() string
}

type Sender interface {
	Send(msg Message, recipient string) (id string, err error)
}

type CachedMessage interface {
	Message
	// Send the message to recipient
	Send(recipient string) (msgid string, err error)
}

type BasicFields struct {
	From    string
	Subject string
	Text    string
	Meta    *MessageMeta
}

// trick for clean struct embedding
type basicFields BasicFields

type MessageMeta struct {
	ContentType string
	Headers     map[string]string
	XSender     string
}
