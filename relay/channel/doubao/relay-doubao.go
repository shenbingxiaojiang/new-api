package doubao

import (
	"fmt"
	relaycommon "one-api/relay/common"
)

func GetRequestURL(info *relaycommon.RelayInfo) (string, error) {
	return fmt.Sprintf("%s/api/v3/chat/completions", info.BaseUrl), nil
}
