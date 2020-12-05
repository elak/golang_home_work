package main

import (
	"bufio"
	"bytes"
	"errors"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

type Environment map[string]EnvValue

// EnvValue helps to distinguish between empty files and files with the first empty line.
type EnvValue struct {
	Value      string
	NeedRemove bool
}

func readFirstLine(dir string, file os.FileInfo) (string, error) {
	path := filepath.Join(dir, file.Name())
	f, err := os.OpenFile(path, os.O_RDONLY, file.Mode().Perm())

	if err != nil {
		return "", err
	}
	defer f.Close()

	reader := bufio.NewReader(f)

	var buffer bytes.Buffer

	for {
		line, isPrefix, err := reader.ReadLine()
		buffer.Write(line)

		// дочитали до конца строки
		if !isPrefix {
			break
		}

		// дочитали до конца файла, не встретив конец строки
		if errors.Is(err, io.EOF) {
			break
		}

		if err != nil {
			return "", err
		}
	}

	// терминальные нули (0x00) заменяются на перевод строки (\n)
	lineStr := string(bytes.ReplaceAll(buffer.Bytes(), []byte{0}, []byte("\n")))
	// пробелы и табуляция в конце удаляются
	lineStr = strings.TrimRight(lineStr, " \t")

	return lineStr, nil
}

// ReadDir reads a specified directory and returns map of env variables.
// Variables represented as files where filename is name of variable, file first line is a value.
func ReadDir(dir string) (Environment, error) {
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	res := make(Environment)

	for _, file := range files {
		valName := file.Name()
		// имя не должно содержать =
		if strings.Contains(valName, "=") {
			continue
		}

		mode := file.Mode()
		if mode&os.ModeSymlink != 0 {
			// пробуем пройти по ссылке
			linkedFile, err := os.Lstat(filepath.Join(dir, valName))
			if err != nil {
				return nil, err
			}
			mode = linkedFile.Mode()
		}

		if !mode.IsRegular() {
			continue
		}

		var val EnvValue

		if file.Size() == 0 {
			val.NeedRemove = true
		} else {
			line, err := readFirstLine(dir, file)
			if err != nil {
				return nil, err
			}
			val.Value = line
		}

		res[valName] = val
	}

	return res, nil
}
