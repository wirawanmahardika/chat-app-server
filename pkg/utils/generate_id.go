package utils

import (
	"fmt"

	gonanoid "github.com/matoous/go-nanoid/v2"
)

func GenerateShortID() string {
	// Menghilangkan 0, O, 1, I, L agar tidak salah ketik
	alphabet := "23456789ABCDEFGHJKLMNPQRSTUVWXYZ"
	id, _ := gonanoid.Generate(alphabet, 8)
	return fmt.Sprintf("VIBE-%s", id)
}
