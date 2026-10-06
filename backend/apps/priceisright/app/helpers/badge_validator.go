package helpers

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

var r *regexp.Regexp

func init() {
	r = regexp.MustCompile("^\\d{1,5}[YACDEFHMPTX]$")
}

func ValidBadgeNumber(badge string) bool {
	letterMap := [...]string{"Y", "A", "C", "D", "E", "F", "H", "M", "P", "T", "X"}
	weights := [...]int{5, 9, 7, 3, 1}

	badgeStr := strings.ToUpper(badge)
	if !r.MatchString(badgeStr) {
		return false
	}

	badgeStr = fmt.Sprintf("%0*s", 5, badge[:len(badge)-1])

	if len(badgeStr) > len(weights) {
		return false
	}

	sum := 0
	for i, numberRune := range badgeStr {
		no, _ := strconv.Atoi(string(numberRune))
		sum += weights[i] * no
	}

	letter := letterMap[sum%len(letterMap)]
	if letter == strings.ToUpper(string(badge[len(badge)-1])) {
		return true
	}

	return false
}
