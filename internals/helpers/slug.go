package helper

import (
	"regexp"
	"strings"
)

func GenerateSlug(input string) string {
	// Ubah ke huruf kecil
	slug := strings.ToLower(input)

	// Ganti semua spasi dan non-alphanumeric jadi strip (-)
	slug = regexp.MustCompile(`[^a-z0-9]+`).ReplaceAllString(slug, "-")

	// Hapus strip di awal/akhir (jika ada)
	slug = strings.Trim(slug, "-")

	return slug
}
