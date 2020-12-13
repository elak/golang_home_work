package hw10_program_optimization //nolint:golint,stylecheck

import (
	"bufio"
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"strings"
)

type User struct {
	ID       int
	Name     string
	Username string
	Email    string
	Phone    string
	Password string
	Address  string
}

var (
	ErrNilInput = errors.New("nil input")
)

type DomainStat map[string]int

func GetDomainStat(r io.Reader, domain string) (DomainStat, error) {
	if r == nil {
		return nil, ErrNilInput
	}

	result := make(DomainStat)

	domain = "." + domain
	domainCheck := domain + `"`

	reader := bufio.NewReader(r)

	var buffer bytes.Buffer
	var user User
	isEOF := false

	for !isEOF {
		line, isPrefix, err := reader.ReadLine()

		if err != nil {
			// дочитали до конца файла, не встретив конец строки?
			if errors.Is(err, io.EOF) {
				isPrefix = false
				isEOF = true
			} else {
				return nil, err
			}
		}

		buffer.Write(line)

		// собрали строку полностью?
		if isPrefix {
			continue
		}

		// Если в строке нет даже намёка на искомый домен - более затратные операции не нужны
		if strings.Contains(buffer.String(), domainCheck) {
			if err := json.Unmarshal(buffer.Bytes(), &user); err == nil {
				// Быстра проверка домена нашла его именно в почте?
				if strings.HasSuffix(user.Email, domain) {
					atPos := strings.Index(user.Email, "@")
					result[strings.ToLower(user.Email[atPos+1:])]++
				}
			}
		}

		buffer.Reset()
	}

	return result, nil
}
