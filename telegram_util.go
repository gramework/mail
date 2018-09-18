package mail

import (
	"crypto/sha256"
	"fmt"
	"strconv"
	"time"
)

func tgGenMsgID(msg Message) string {
	return fmt.Sprintf(
		"%x",
		sha256.Sum256(append([]byte(msg.Source()), strconv.FormatInt(time.Now().UnixNano(), 10)...)),
	)
}
