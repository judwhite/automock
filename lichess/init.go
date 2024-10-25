package lichess

import "os"

var lichessAPIToken string

func init() {
	lichessAPIToken = os.Getenv("LICHESS_API_TOKEN")
	if lichessAPIToken == "" {
		panic("LICHESS_API_TOKEN not set")
	}
}
