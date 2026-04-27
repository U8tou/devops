package devflow

import (
	"os"
	"strings"
)

// OsEnvWithProcess 将流程环境变量覆盖到当前进程环境之后（同名键以 proc 为准）。
func OsEnvWithProcess(base []string, proc map[string]string) []string {
	if len(proc) == 0 {
		return base
	}
	m := envPairsToMap(base)
	for k, v := range proc {
		m[k] = v
	}
	return mapToEnvSlice(m)
}

func envPairsToMap(env []string) map[string]string {
	m := make(map[string]string)
	for _, e := range env {
		idx := strings.IndexByte(e, '=')
		if idx <= 0 {
			continue
		}
		m[e[:idx]] = e[idx+1:]
	}
	return m
}

func mapToEnvSlice(m map[string]string) []string {
	out := make([]string, 0, len(m))
	for k, v := range m {
		out = append(out, k+"="+v)
	}
	return out
}

func defaultEnviron() []string {
	return os.Environ()
}
