package service

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
