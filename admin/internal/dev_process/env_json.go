package dev_process

import (
	"encoding/json"
	"fmt"
	"strings"
)

func envJSONToMap(s string) map[string]string {
	out := make(map[string]string)
	s = strings.TrimSpace(s)
	if s == "" {
		return out
	}
	var raw map[string]interface{}
	if err := json.Unmarshal([]byte(s), &raw); err != nil {
		return out
	}
	for k, v := range raw {
		out[k] = fmt.Sprint(v)
	}
	return out
}

func mapToEnvJSON(m map[string]string) (string, error) {
	if m == nil {
		return "{}", nil
	}
	for k := range m {
		if strings.TrimSpace(k) == "" {
			return "", fmt.Errorf("环境变量名不能为空")
		}
	}
	b, err := json.Marshal(m)
	if err != nil {
		return "", err
	}
	return string(b), nil
}
