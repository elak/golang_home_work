package hw03_frequency_analysis //nolint:golint,stylecheck

import (
	"regexp"
	"sort"
	"strings"
)

func Top10(text string) []string {
	wordsCount := map[string]int{}

	// Словом считается набор символов, разделенных пробельными символами и знаками припинания
	// в терминах класов юникода "символами" будем считать буквы L, цифры N, символы S и дефис
	re := regexp.MustCompile(`([-\pL\pN\pS]+)\b*`)
	words := re.FindAllString(text, -1)

	for _, word := range words {
		if word == "-" {
			continue // "(тире) - это не слово", и дефис - тоже
		}

		word = strings.ToLower(word) // "Нога" и "нога" - это одинаковые слова

		_, exist := wordsCount[word]
		if exist {
			wordsCount[word]++
		} else {
			wordsCount[word] = 1
		}
	}

	keys := make([]string, 0, len(wordsCount))
	for k := range wordsCount {
		keys = append(keys, k)
	}

	sort.Slice(keys, func(i, j int) bool {
		return wordsCount[keys[i]] > wordsCount[keys[j]]
	})

	upperIdx := len(keys)
	if upperIdx > 10 {
		upperIdx = 10
	}

	return keys[0:upperIdx]
}
