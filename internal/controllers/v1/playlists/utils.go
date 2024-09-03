package playlists

import "math/rand"

var letters = []rune("ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")

func GeneratePlaylistCode(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}
