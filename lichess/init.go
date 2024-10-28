package lichess

import (
	"fmt"
	"net/http"
	"os"
)

var lichessAPIToken string
var authHeader = make(http.Header)

func init() {
	lichessAPIToken = os.Getenv("LICHESS_API_TOKEN")
	if lichessAPIToken == "" {
		panic("LICHESS_API_TOKEN not set")
	}

	authHeader.Set("Authorization", fmt.Sprintf("Bearer %s", lichessAPIToken))
}
