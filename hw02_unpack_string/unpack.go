package hw02unpackstring

import (
	"errors"
	"strconv"
	"strings"
	"unicode"
)

var ErrInvalidString = errors.New("invalid string")

func checkErrors(inp string) bool {
	runes := []rune(inp)

	// Check condition if input string is a number
	if _, errChk := strconv.Atoi(inp); errChk == nil {
		return true
	}

	// Check condition if first char is a number
	for idx, char := range inp {
		if unicode.IsDigit(char) && idx == 0 {
			return true
		}
	}

	for i := 2; i < len(runes); i++ {
		if unicode.IsDigit(runes[i]) && unicode.IsDigit(runes[i-1]) && (i > 1 && runes[i-2] != 92) {
			return true
		}
	}
	return false
}

func repeater(rep int, str string, run rune) string {
	var strBldr strings.Builder
	if rep > 0 {
		strBldr.WriteString(str)
		strBldr.WriteString(strings.Repeat(string(run), rep))
		str = strBldr.String()
	} else {
		str = strings.Replace(str, string(run), "", 1)
	}
	return str
}

func Unpack(inp string) (string, error) {
	var strBldr strings.Builder
	runes := []rune(inp)

	if checkErrors(inp) {
		return "", ErrInvalidString
	}

	for i := 0; i < len(runes); i++ {
		// Check for backslas

		switch rn := runes[i]; {
		case rn == 92:
			i++
			strBldr.WriteRune(runes[i])
		case unicode.IsDigit(rn):
			rep, _ := strconv.Atoi(string(runes[i]))
			rep--

			tempStr := strBldr.String()
			tempStr = repeater(rep, tempStr, runes[i-1])
			strBldr.Reset()
			strBldr.WriteString(tempStr)
		default:
			strBldr.WriteRune(runes[i])
		}
	}
	return strBldr.String(), nil
}
