package mail_test

import (
	"testing"

	"github.com/gramework/mail"
)

func TestPlain(t *testing.T) {
	pt := mail.NewPlain(&mail.BasicFields{
		Text: "",
	})

	_ = pt
}
