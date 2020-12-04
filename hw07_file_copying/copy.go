package main

import (
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

var (
	ErrUnsupportedFile       = errors.New("unsupported file")
	ErrOffsetExceedsFileSize = errors.New("offset exceeds file size")
	ErrSameFile              = errors.New("can not copy file into itself")
	ErrIsDir                 = errors.New("directories are not supported")
)

type progressBar struct {
	maxDisplayLength     int64  // максимальная длина шкалы в символах
	currentDisplayLength int64  // текущая длина шкалы в символах
	goalProgress         int64  // максимальное значение програсса
	currentProgress      int64  // текущее значение програсса
	displayChunk         string // отображение одного деления шкалы
}

func (thisBar *progressBar) progress(amount int64) {
	if thisBar.currentProgress == thisBar.goalProgress {
		return
	}

	thisBar.currentProgress += amount
	progressLength := thisBar.maxDisplayLength * thisBar.currentProgress / thisBar.goalProgress
	for progressLength > thisBar.currentDisplayLength {
		thisBar.currentDisplayLength++
		fmt.Print(thisBar.displayChunk)
	}

	if thisBar.currentProgress == thisBar.goalProgress {
		fmt.Println()
	}
}

func (thisBar *progressBar) init(goal int64, length int64) {
	thisBar.maxDisplayLength = length
	thisBar.currentDisplayLength = 0
	thisBar.goalProgress = goal
	thisBar.displayChunk = "."

	if goal == 0 {
		return
	}

	if length > goal {
		thisBar.displayChunk = strings.Repeat(thisBar.displayChunk, int(length/goal))
		thisBar.maxDisplayLength = goal
	}
}

func checkPaths(fromPath, toPath *string) error {
	fromPathAbs, err := filepath.Abs(*fromPath)
	if err != nil {
		return err
	}

	toPathAbs, err := filepath.Abs(*toPath)
	if err != nil {
		return err
	}

	if fromPathAbs == toPathAbs {
		return ErrSameFile
	}

	*fromPath = fromPathAbs
	*toPath = toPathAbs

	return nil
}

func checkFileStat(stat os.FileInfo) error {
	mode := stat.Mode()

	if mode&os.ModeSymlink == os.ModeSymlink {
		// пробуем пройти по ссылке
		linkStat, err := os.Lstat(stat.Name())
		if err != nil {
			return err
		}

		return checkFileStat(linkStat)
	}

	// если файл существует - это не должен быть каталог
	if mode.IsDir() {
		return ErrIsDir
	}

	// ... и вообще что-либо кроме файла
	if !mode.IsRegular() {
		return ErrUnsupportedFile
	}

	return nil
}

func tryOpenFile(name string, flag int, perm os.FileMode) (*os.File, error) {
	fStat, err := os.Stat(name)

	if err == nil {
		perm = fStat.Mode().Perm()
		err = checkFileStat(fStat)
	}

	if os.IsNotExist(err) {
		// если файла нет то мы должны быть в режиме создания
		if flag&os.O_CREATE == 0 {
			return nil, err
		}
	} else if err != nil {
		return nil, err
	}

	f, err := os.OpenFile(name, flag, perm)

	if err != nil {
		return nil, err
	}

	return f, nil
}

func copyData(from io.Reader, to io.Writer, amount int64) error {
	if amount == 0 {
		return nil
	}

	var pb progressBar
	pb.init(amount, 40)
	chunkSize := amount / 40
	if chunkSize == 0 {
		chunkSize = amount
	}

	for {
		written, err := io.CopyN(to, from, chunkSize)
		pb.progress(written)
		amount -= written

		if amount == 0 {
			break
		}

		if errors.Is(err, io.EOF) {
			return io.ErrUnexpectedEOF
		}

		if err != nil {
			return err
		}
	}

	return nil
}

func Copy(fromPath, toPath string, offset, limit int64) error {
	// преобразуем пути в абсолютные и проверим не станут ли они совпадать
	err := checkPaths(&fromPath, &toPath)
	if err != nil {
		return err
	}

	// попробуем открыть входной файл и проверим ограничения параметров
	fSrc, err := tryOpenFile(fromPath, os.O_RDONLY, 0666)

	if err != nil {
		return err
	}

	defer fSrc.Close()

	fStat, err := fSrc.Stat()
	if err != nil {
		return err
	}

	// * offset больше, чем размер файла - невалидная ситуация;
	if offset > fStat.Size() {
		return ErrOffsetExceedsFileSize
	}

	_, err = fSrc.Seek(offset, 0)
	if err != nil {
		return err
	}

	// * limit больше, чем размер файла - валидная ситуация, копируется исходный файл до его EOF;
	if limit == 0 || offset+limit > fStat.Size() {
		limit = fStat.Size() - offset
	}

	// попробуем открыть или создать выходной файл
	fDst, err := tryOpenFile(toPath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, fStat.Mode().Perm())

	if err != nil {
		return err
	}

	defer fDst.Close()

	// и наконец собственно переложим байтики из одного файла в другой
	err = copyData(fSrc, fDst, limit)

	if err != nil {
		return err
	}

	return nil
}
