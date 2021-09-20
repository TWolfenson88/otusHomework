package hw02unpackstring

import (
	"errors"
	"strconv"
	"strings"
	"unicode"
)

var ErrInvalidString = errors.New("invalid string")

func Unpack(inp string) (string, error) {
	var strBldr strings.Builder
	runes := []rune(inp)
	for i := 0; i < len(runes); i++ {
		//Check for backslas
		if runes[i] == 92 {
			i++
			strBldr.WriteRune(runes[i])
		} else if rep, err := strconv.Atoi(string(runes[i])); err == nil {
			//Check for non-digit input
			if (i > 0 && unicode.IsDigit(runes[i-1]) && (i > 1 && runes[i-2] != 92)) || i == 0 {
				return "", ErrInvalidString
			}
			//Check for zero digit
			rep = rep - 1
			if rep > 0 {
				strBldr.WriteString(strings.Repeat(string(runes[i-1]), rep))
			} else {
				tempStr := strBldr.String()
				tempStr = strings.Replace(tempStr, string(runes[i-1]), "", 1)
				strBldr.Reset()
				strBldr.WriteString(tempStr)
			}

		} else {
			strBldr.WriteRune(runes[i])
		}
	}
	return strBldr.String(), nil
}
