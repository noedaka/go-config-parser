package service

import (
	"fmt"
	"strings"
)

type severity string

const (
	Low    severity = "LOW"
	Medium severity = "MEDIUM"
	High   severity = "HIGH"
)

type Issue struct {
	Severity       severity
	Message        string
	Recommendation string
}

type Rule interface {
	Name() string
	Check(data any) []Issue
}

func FormatIssues(issues []Issue) string {
	var sb strings.Builder
	for _, issue := range issues {
		fmt.Fprintf(&sb, "[%s] %s\nРекомендация: %s\n\n",
			issue.Severity, issue.Message, issue.Recommendation)
	}
	return sb.String()
}
