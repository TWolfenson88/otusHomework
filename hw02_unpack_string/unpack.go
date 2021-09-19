package hw02unpackstring

import (
	"errors"
	"strconv"
	"strings"
	"unicode"
)

var ErrInvalidString = errors.New("invalid string")

func Unpack(inp string) (string, error) {

	//Check condition if input string is a number
	_, errChk := strconv.Atoi(inp)
	if errChk == nil {
		return "", ErrInvalidString
	}

	//Check condition if first char is a number
	for idx, char := range inp {
		if unicode.IsDigit(char) && idx == 0 {
			return "", ErrInvalidString
		}
	}

	//Check condition if input sting contian not a digit
	prevChar := 'a'
	for _, char := range inp {
		if unicode.IsDigit(char) && unicode.IsDigit(prevChar) {
			return "", ErrInvalidString
		}
		prevChar = char
	}

	outStr := ""

	// charr := ' '
	prevChar = ' '

	//a4bc5      aab5e

	for _, char := range inp {

		if unicode.IsLetter(prevChar) && unicode.IsDigit(char) {
			rep, err := strconv.Atoi(string(char))
			if err != nil {
				return "", ErrInvalidString
			}
			rep = rep - 1
			if rep > 0 {
				outStr += strings.Repeat(string(prevChar), rep)
			} else if rep < 0 {
				outStr = strings.Replace(outStr, string(prevChar), "", 1)
			}
			prevChar = char
		} else if unicode.IsDigit(prevChar) && unicode.IsLetter(char) {
			outStr += string(char)
			prevChar = char
		} else if unicode.IsLetter(prevChar) && unicode.IsLetter(char) {
			outStr += string(char)
			prevChar = char
		} else if prevChar == ' ' {
			outStr += string(char)
			prevChar = char
		}
	}

	return outStr, nil
}
