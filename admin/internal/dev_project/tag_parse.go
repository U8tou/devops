package dev_project

import (
	"strconv"
	"strings"
)

func parseCommaInt64s(s string) []int64 {
	s = strings.TrimSpace(s)
	if s == "" {
		return nil
	}
	parts := strings.Split(s, ",")
	var out []int64
	for _, p := range parts {
		p = strings.TrimSpace(p)
		if p == "" {
			continue
		}
		n, err := strconv.ParseInt(p, 10, 64)
		if err != nil {
			continue
		}
		out = append(out, n)
	}
	return out
}

func parseTagOther(s string) bool {
	s = strings.TrimSpace(strings.ToLower(s))
	return s == "1" || s == "true" || s == "yes"
}

func parseTagModeExclude(s string) bool {
	s = strings.TrimSpace(strings.ToLower(s))
	return s == "exclude"
}
