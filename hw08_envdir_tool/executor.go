package main

import (
	"errors"
	"os"
	"os/exec"
	"strings"
)

func patchEnvironment(envPatch Environment) []string {
	if len(envPatch) == 0 {
		return nil
	}
	// для переопределения значений достаточно было бы дописать наши в конец
	// но удалить что-то кроме как обходом с копированием не получится
	patchedEnv := make([]string, 0, len(os.Environ())+len(envPatch))

	for _, kv := range os.Environ() {
		eq := strings.Index(kv, "=")
		if eq < 0 {
			patchedEnv = append(patchedEnv, kv)
			continue
		}
		k := kv[:eq]

		patchVal, usePatch := envPatch[k]

		if usePatch {
			if patchVal.NeedRemove {
				continue
			}
			kv = k + "=" + patchVal.Value
			delete(envPatch, k)
		}

		patchedEnv = append(patchedEnv, kv)
	}

	for k, v := range envPatch {
		if v.NeedRemove {
			continue
		}
		patchedEnv = append(patchedEnv, k+"="+v.Value)
	}

	return patchedEnv
}

// RunCmd runs a command + arguments (cmd) with environment variables from env.
func RunCmd(cmd []string, env Environment) (returnCode int) {
	cmdWrap := exec.Command(cmd[0], cmd[1:]...) //nolint:gosec

	// стандартные потоки ввода/вывода/ошибок пробрасывались в вызываемую программу;
	cmdWrap.Stderr = os.Stderr
	cmdWrap.Stdin = os.Stdin
	cmdWrap.Stdout = os.Stdout

	cmdWrap.Env = patchEnvironment(env)

	err := cmdWrap.Run()

	// код выхода утилиты совпадал с кодом выхода программы.
	if err != nil {
		returnCode = -1

		var exitErr *exec.ExitError
		if errors.As(err, &exitErr) {
			returnCode = exitErr.ExitCode()
		}
	}

	return
}
