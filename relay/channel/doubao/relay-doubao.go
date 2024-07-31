package doubao

import (
	"fmt"
	relaycommon "one-api/relay/common"
)

func GetRequestURL(info *relaycommon.RelayInfo) (string, error) {
	// Support BotChatCompletions
	if info.RequestModelName == "Doubao-bot-chat" {
		return fmt.Sprintf("%s/api/v3/bots/chat/completions", info.BaseUrl), nil
	}
	return fmt.Sprintf("%s/api/v3/chat/completions", info.BaseUrl), nil
}
