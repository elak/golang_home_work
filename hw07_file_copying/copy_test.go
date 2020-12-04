package main

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCopy(t *testing.T) {
	t.Run("directories", func(t *testing.T) {
		err := Copy("./testdata", "out.txt", 0, 0)
		require.Equal(t, ErrIsDir, err)

		err = Copy("testdata/input.txt", "./testdata", 0, 0)
		require.Equal(t, ErrIsDir, err)
	})

	t.Run("unsupported", func(t *testing.T) {
		err := Copy("/dev/urandom", "out.txt", 0, 0)
		require.Equal(t, ErrUnsupportedFile, err)

		err = Copy("testdata/input.txt", "/dev/urandom", 0, 0)
		require.Equal(t, ErrUnsupportedFile, err)
	})

	t.Run("self copy", func(t *testing.T) {
		err := Copy("testdata/input.txt", "testdata/input.txt", 0, 0)
		require.Equal(t, ErrSameFile, err)

		err = Copy("testdata/input.txt", "testdata/../testdata/input.txt", 0, 0)
		require.Equal(t, ErrSameFile, err)
	})

	t.Run("wrong offset", func(t *testing.T) {
		err := Copy("testdata/input.txt", "out.txt", 100000, 0)
		require.Equal(t, ErrOffsetExceedsFileSize, err)
	})

}
