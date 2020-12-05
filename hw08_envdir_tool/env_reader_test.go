package main

import (
	"github.com/stretchr/testify/require"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestReadDir(t *testing.T) {
	fullPath, _ := filepath.Abs("testdata/empty")

	func() {
		// создаём пустой каталог
		err := os.Mkdir(fullPath, 0666)
		require.Nil(t, err)
		require.DirExists(t, fullPath)
	}()

	t.Cleanup(func() {
		// удаляем пустой каталог
		_ = os.Remove(fullPath)
	})

	t.Run("presets", func(t *testing.T) {
		env, err := ReadDir("testdata/env")
		require.Nil(t, err)

		require.Equal(t, 5, len(env))

		// правильный замена 0х00
		require.True(t, strings.Contains(env["FOO"].Value, "\n"))
		// Игнорирование второй строки
		require.False(t, strings.Contains(env["BAR"].Value, "\n"))

		// правильный разбор
		require.True(t, env["UNSET"].NeedRemove)
	})

	// пустые данные
	t.Run("presets", func(t *testing.T) {
		env, err := ReadDir("testdata/empty")
		require.Nil(t, err)
		require.Empty(t, env)

	})

	// правильный разбор
	t.Run("skip dir", func(t *testing.T) {
		env, err := ReadDir("testdata")
		require.Nil(t, err)
		// не считывает директории
		require.Equal(t, 1, len(env))
		// считал первую строку крипта
		require.Equal(t, "#!/usr/bin/env bash", env["echo.sh"].Value)
	})
}
