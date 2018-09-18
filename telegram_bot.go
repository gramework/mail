package mail

import (
	"strings"
	"sync"

	"github.com/google/uuid"
	"gitlab.com/toby3d/telegram"
)

type tgEntity struct {
	sync.RWMutex
	chans map[string][]int64
}

type tgReg map[string]*tgEntity

type tgMultibotSubscribers struct {
	sync.RWMutex
	bots tgReg
}

var tgRegister = tgMultibotSubscribers{
	bots: make(tgReg),
}

const tgWelcomeMsg = `Hello.

To create a channel and get a new channelUUID, use /create.
Note: You'll be autosubscribed to your new channel.

To subscribe to a channel, use /subscribe channelUUID

` + FullVersion

func (tg *TelegramSubSender) serveBot() {
	// we should allow only message update
	updates := tg.bot.NewLongPollingChannel(&telegram.GetUpdatesParameters{
		AllowedUpdates: []string{"message"},
	})
	// TODO: consider a webhook-based sender bot.NewWebhookChannel() and tg login

	for update := range updates {
		// bullet-proof double check if it is a message
		if !update.IsMessage() {
			// not a message
			tg.log.
				WithField("got", update.Type()).
				Warn("got non-message update with filtered message-only request")
			continue
		}

		// Ok, it's definitely a message. From here on use
		// only this struct to avoid unnecessary panics
		u := update.Message

		go tg.processUpdate(u)
	}
}

func (tg *TelegramSubSender) warmup() {
	tg.bot, tg.botInitErr = telegram.New(tg.AccessToken)
	if tg.botInitErr == nil {
		tg.log = log.WithField("submod", "TelegramSubSender")
		// accessToken can be mutated or removed after initialization
		// but we should track the same bot anyway
		tg.botKey = tg.AccessToken
		tgRegister.Lock()
		if _, ok := tgRegister.bots[tg.botKey]; !ok {
			tgRegister.bots[tg.botKey] = &tgEntity{
				chans: make(map[string][]int64),
			}
		}
		tgRegister.Unlock()
		go tg.serveBot()
	}
}

func (tg *TelegramSubSender) Warmup() error {
	if len(tg.AccessToken) == 0 {
		NewCredReqErr("AccessToken", "TelegramSub")
	}
	tg.initBot.Do(tg.warmup)
	return tg.botInitErr
}

func (tg *TelegramSubSender) sendMsg(msg *telegram.SendMessageParameters) (err error) {
	for i := 0; i < 8; i++ {
		_, err = tg.bot.SendMessage(msg)
		if err == nil {
			return nil
		}
	}
	return err
}

func (tg *TelegramSubSender) processUpdate(u *telegram.Message) {
	// we work only in private chats
	if !u.Chat.IsPrivate() {
		tg.sendMsg(telegram.NewMessage(u.Chat.ID, "OH MY... So much humans!"))
		tg.sendMsg(telegram.NewMessage(u.Chat.ID, "/me scared"))
		tg.bot.LeaveChat(u.Chat.ID)
		tg.log.Warn("not a private chat")
		return
	}

	// check if this command send to the bot
	if !tg.bot.IsCommandToMe(u) {
		// nope, this msg seems to be sent to anyone else
		return
	}

	switch {
	case u.IsCommandEqual("start"):
		tg.sendMsg(telegram.NewMessage(u.Chat.ID, tgWelcomeMsg))
	case u.IsCommandEqual("create"):
		tgRegister.RLock()
		botReg := tgRegister.bots[tg.botKey]
		tgRegister.RUnlock()

		botReg.Lock()
		for {
			chanUUID := uuid.New().String()
			if _, ok := botReg.chans[chanUUID]; !ok {
				botReg.chans[chanUUID] = []int64{
					u.Chat.ID,
				}
				tg.sendMsg(telegram.NewMessage(
					u.Chat.ID,
					"Ok, we've got you covered. Here's your new channel UUID: "+chanUUID,
				))
				break
			}
		}
		botReg.Unlock()

	case u.IsCommandEqual("subscribe"):
		uuid := u.CommandArgument()
		if len(strings.TrimSpace(uuid)) == 0 {
			tg.sendMsg(telegram.NewMessage(
				u.Chat.ID, "Subscribe to... what?",
			))
			return
		}

		tgRegister.RLock()
		botReg := tgRegister.bots[tg.botKey]
		tgRegister.RUnlock()
		botReg.Lock()
		if _, ok := botReg.chans[uuid]; !ok {
			tg.sendMsg(telegram.NewMessage(
				u.Chat.ID,
				"Wait, I don't have such channel",
			))
			botReg.Unlock()
			return
		}

		for _, v := range botReg.chans[uuid] {
			if v == u.Chat.ID {
				tg.sendMsg(telegram.NewMessage(
					u.Chat.ID,
					"You already subscribed --.",
				))
				botReg.Unlock()
				return
			}
		}

		botReg.chans[uuid] = append(botReg.chans[uuid], u.Chat.ID)

		botReg.Unlock()
	}
}
