package traceid

import (
	"math/rand"
	"time"
)

func ID() string {
	rand.New(rand.NewSource(time.Now().UnixNano()))

	chars := []rune("ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789")
	id := make([]rune, 12)
	for i := range id {
		id[i] = chars[rand.Intn(len(chars))]
	}

	return string(id)
}
