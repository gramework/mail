package mail

type PlainText struct {
	*basicFields
}

func NewPlain(fields *BasicFields) *PlainText {
	return &PlainText{
		basicFields: (*basicFields)(fields),
	}
}

func (p *PlainText) Source() string {
	return p.basicFields.Text
}

func (p *PlainText) Meta() *MessageMeta {
	return p.basicFields.Meta
}

func (p *PlainText) Fields() *BasicFields {
	return (*BasicFields)(p.basicFields)
}
