package hw02_unpack_string //nolint:golint,stylecheck

import (
	"errors"
	"log"
	"strconv"
	"strings"
	"unicode"
)

var ErrInvalidString = errors.New("invalid string")

func Unpack(packedStr string) (string, error) {
	escChar := false // текущий символ экранирован предыдущим
	skipChar := true // символ из буфера не надо выводить в результат

	var prevChar rune = -1 // буферный символ

	var result strings.Builder

	for _, curChar := range packedStr {
		if !escChar && prevChar == '\\' {
			// неэкранированный обратный слэш
			// взводим флаг и идём дальше, не дополняя результат
			escChar = true
			prevChar = curChar

			continue
		}

		escChar = false

		if unicode.IsDigit(curChar) {
			if skipChar {
				// мы или в самом начале строки или сразу после распаковки пары символ-цифра
				// т.е. это или стока, начинающаяся с цифры, или вторая цифра подряд,
				// а это - полный перечень возможных ошибок формата
				return "", ErrInvalidString
			}

			chunkLength, err := strconv.Atoi(string(curChar))

			if err != nil {
				log.Fatalf("Error : %v\n", err)
			}

			chunk := strings.Repeat(string(prevChar), chunkLength)
			result.WriteString(chunk)

			skipChar = true

			continue
		}

		if !skipChar {
			result.WriteRune(prevChar)
		}

		skipChar = false
		prevChar = curChar
	}

	if !skipChar {
		if !escChar && prevChar == '\\' {
			// строку закончившуюся на неэкранированный обратный слэш считаем ошибкой формата
			return "", ErrInvalidString
		}

		result.WriteRune(prevChar)
	}

	return result.String(), nil
}
