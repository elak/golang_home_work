package logger

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestLogger(t *testing.T) {
	logsPath, _ := filepath.Abs("testdata")

	var loger *Logger
	var err error

	logMessages := func(lvl string) string {
		logPath := filepath.Join(logsPath, lvl+".log")

		loger, err = New(lvl, logPath)
		require.NoError(t, err)

		loger.Error("ERROR message")
		loger.Warning("WARNING message")
		loger.Info("INFO message")
		loger.Debug("DEBUG message")

		err = loger.Close()
		require.NoError(t, err)

		logFile, err := os.Open(logPath)
		require.NoError(t, err)

		logText, err := ioutil.ReadAll(logFile)
		require.NoError(t, err)

		defer func() {
			err = logFile.Close()
			require.NoError(t, err)

			err = os.Remove(logPath)
			require.NoError(t, err)
		}()

		return string(logText)
	}

	msg := ""
	msg = logMessages("ERROR")
	require.NotEqual(t, -1, strings.Index(msg, "ERROR"))
	require.Equal(t, -1, strings.Index(msg, "WARNING"))
	require.Equal(t, -1, strings.Index(msg, "INFO"))
	require.Equal(t, -1, strings.Index(msg, "DEBUG"))

	msg = logMessages("WARNING")
	require.NotEqual(t, -1, strings.Index(msg, "ERROR"))
	require.NotEqual(t, -1, strings.Index(msg, "WARNING"))
	require.Equal(t, -1, strings.Index(msg, "INFO"))
	require.Equal(t, -1, strings.Index(msg, "DEBUG"))

	msg = logMessages("INFO")
	require.NotEqual(t, -1, strings.Index(msg, "ERROR"))
	require.NotEqual(t, -1, strings.Index(msg, "WARNING"))
	require.NotEqual(t, -1, strings.Index(msg, "INFO"))
	require.Equal(t, -1, strings.Index(msg, "DEBUG"))

	msg = logMessages("DEBUG")
	require.NotEqual(t, -1, strings.Index(msg, "ERROR"))
	require.NotEqual(t, -1, strings.Index(msg, "WARNING"))
	require.NotEqual(t, -1, strings.Index(msg, "INFO"))
	require.NotEqual(t, -1, strings.Index(msg, "DEBUG"))
}
