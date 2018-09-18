package mail

import "github.com/pkg/errors"

func NewCredReqErr(fieldName, senderType string) error {
	return errors.Errorf("field required but was not provided %20s    %15s", "fieldName="+fieldName, "senderType="+senderType)
}

func NewNoSuchRecipientErr(recipient, senderType string) error {
	return errors.Errorf("no such recipient %20s    %15s", "recipient="+recipient, "senderType="+senderType)
}
