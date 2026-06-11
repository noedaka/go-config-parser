package rules

import (
	"strings"

	"github.com/noedaka/go-config-parser/internal/service"
)

type DebugLogRule struct{}

func (r DebugLogRule) Name() string { return "debug-logging" }

func (r DebugLogRule) Check(data any) []service.Issue {
	var issues []service.Issue
	walk(data, func(key string, value any) {
		if strings.EqualFold(key, "level") || strings.EqualFold(key, "log_level") {
			if s, ok := value.(string); ok && strings.EqualFold(s, "debug") {
				issues = append(issues, service.Issue{
					Severity:       service.Low,
					Message:        "логирование в debug-режиме",
					Recommendation: "Поменяйте режим на более избирательный (info и выше)",
				})
			}
		}
	})
	return issues
}
