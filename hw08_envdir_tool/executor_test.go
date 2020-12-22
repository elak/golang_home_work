package main

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestRunCmd(t *testing.T) {

	t.Run("simple", func(t *testing.T) {
		args := []string{"echo", "1"}

		retCode := RunCmd(args, make(Environment))
		require.Equal(t, 0, retCode)
	})

	t.Run("no args", func(t *testing.T) {
		args := []string{"echo"}

		retCode := RunCmd(args, make(Environment))
		require.Equal(t, 0, retCode)
	})

	t.Run("no cmd", func(t *testing.T) {
		args := []string{""}

		retCode := RunCmd(args, make(Environment))
		require.Equal(t, -1, retCode)
	})
}
