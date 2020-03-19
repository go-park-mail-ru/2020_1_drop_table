package owners

import (
	_ "github.com/lib/pq"
	redisStore "gopkg.in/boj/redistore.v1"
	"os"
)

var CookieStore, err = redisStore.NewRediStore(10, "tcp", ":6379",
	"", []byte(os.Getenv("SESSION_KEY")))

func init() {
	if err != nil {
		panic("Can't connect to redis")
	}
}

const CookieName = "authCookie"
