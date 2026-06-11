package service

import "strings"

type DebugLogRule struct{}

func (r DebugLogRule) Name() string { return "debug-logging" }

func (r DebugLogRule) Check(data any) []Issue {
	var issues []Issue
	walk(data, func(key string, value any) {
		if strings.EqualFold(key, "level") || strings.EqualFold(key, "log_level") {
			if s, ok := value.(string); ok && strings.EqualFold(s, "debug") {
				issues = append(issues, Issue{
					Severity:       Low,
					Message:        "логирование в debug-режиме",
					Recommendation: "Поменяйте режим на более избирательный (info и выше)",
				})
			}
		}
	})
	return issues
}
