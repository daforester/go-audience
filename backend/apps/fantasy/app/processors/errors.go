package processors

import (
	"errors"
	"github.com/spf13/viper"
	"strings"
	"unicode"
	"unicode/utf8"
)

var (
	ErrGameNotFound   = errors.New("game not found")
	ErrPlayerNotFound = errors.New("player not found")
	ErrForbidden      = errors.New("access denied")
	ErrInvalid        = errors.New("invalid request")
	ErrWrongPhase     = errors.New("not allowed at this point in the game")
	ErrTooSoon        = errors.New("slow down, try again in a moment")
)

// CleanText trims whitespace, strips control characters and limits the length in runes
func CleanText(s string, maxLength int) string {
	s = strings.Map(func(r rune) rune {
		if r == '\n' || r == '\t' {
			return ' '
		}
		if unicode.IsControl(r) || r == utf8.RuneError {
			return -1
		}
		return r
	}, s)
	s = strings.Join(strings.Fields(s), " ")

	if maxLength > 0 && utf8.RuneCountInString(s) > maxLength {
		s = string([]rune(s)[:maxLength])
	}

	return s
}

func CleanName(name string) string {
	return CleanText(name, viper.GetInt("MAX_NAME_LENGTH"))
}
