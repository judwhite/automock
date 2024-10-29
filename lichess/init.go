package lichess

import (
	"fmt"
	"net/http"
	"os"

	"automock/utils"
)

var lichessAPIToken string
var authHeader = make(http.Header)

const apiTokenEnvName = "LICHESS_API_TOKEN"

func init() {
	lichessAPIToken = os.Getenv(apiTokenEnvName)
	if lichessAPIToken == "" {
		utils.Log(fmt.Sprintf("NOTE: Set `%s` env for a slightly better Lichess API experience.", apiTokenEnvName))
		return
	}

	authHeader.Set("Authorization", fmt.Sprintf("Bearer %s", lichessAPIToken))
}
