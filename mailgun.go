package mail

import mailgun "github.com/mailgun/mailgun-go"

type MailgunSender struct {
	Domain        string
	APIKey        string
	ValidationKey string
}

var _ Sender = &MailgunSender{}

func (mg *MailgunSender) Send(msg Message, recipient string) (id string, err error) {
	switch 0 {
	case len(mg.Domain):
		return "", NewCredReqErr("Domain", "mailgun")
	case len(mg.APIKey):
		return "", NewCredReqErr("APIKey", "mailgun")
	case len(mg.ValidationKey):
		return "", NewCredReqErr("ValidationKey", "mailgun")
	}

	mgclient := mailgun.NewMailgun(mg.Domain, mg.APIKey, mg.ValidationKey)

	f := msg.Fields()
	mgmsg := mgclient.NewMessage(f.From, f.Subject, f.Text, recipient)
	meta := msg.Meta()
	for name, value := range meta.Headers {
		mgmsg.AddHeader(name, value)
	}
	xsender := FullVersion
	if len(meta.XSender) > 0 {
		xsender = meta.XSender
	}
	mgmsg.AddHeader("X-Sender", xsender)
	_, id, err = mgclient.Send(mgmsg)
	return id, err
}
