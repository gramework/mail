package mail

import (
	"sync"

	alog "github.com/apex/log"

	"gitlab.com/toby3d/telegram"
)

type TelegramSubSender struct {
	AccessToken   string
	SendDebugInfo bool
	botKey        string
	initBot       sync.Once
	log           *alog.Entry
	bot           *telegram.Bot
	botInitErr    error
}

func (tg *TelegramSubSender) Send(msg Message, recipient string) (id string, err error) {
	if len(tg.AccessToken) == 0 {
		return "", NewCredReqErr("AccessToken", "TelegramSub")
	}
	tg.initBot.Do(tg.warmup)
	msgID := tgGenMsgID(msg)

	if tg.botInitErr != nil {
		tg.initBot = sync.Once{}
		return "", tg.botInitErr
	}

	tgRegister.RLock()
	botReg := tgRegister.bots[tg.botKey]
	tgRegister.RUnlock()

	botReg.RLock()
	subs, ok := botReg.chans[recipient]
	botReg.RUnlock()
	if !ok {
		return "", NewNoSuchRecipientErr(recipient, "TelegramSub")
	}
	msgSrc := msg.Source()
	if tg.SendDebugInfo {
		msgSrc += "\n\n"
		msgSrc += "**UUID:** " + msgID
	}
	for _, chatID := range subs {
		tgMsg := &telegram.SendMessageParameters{
			ParseMode: "Markdown",
			ChatID:    chatID,
			Text:      msgSrc,
		}
		tg.sendMsg(tgMsg)
	}

	return msgID, nil
}
