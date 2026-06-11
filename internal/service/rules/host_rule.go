package rules

import (
	"fmt"
	"strings"

	"github.com/noedaka/go-config-parser/internal/service"
)

type ZeroHostRule struct{}

func (r ZeroHostRule) Name() string { return "bind-all-interfaces" }

func (r ZeroHostRule) Check(data any) []service.Issue {
	var issues []service.Issue
	hostKeys := []string{"host", "bind", "listen", "address"}
	walk(data, func(key string, value any) {
		k := strings.ToLower(key)
		for _, hk := range hostKeys {
			if k == hk {
				if s, ok := value.(string); ok && (s == "0.0.0.0" || s == "[::]") {
					issues = append(issues, service.Issue{
						Severity:       service.Medium,
						Message:        fmt.Sprintf("сервис слушает на %s (все интерфейсы) без видимых ограничений", s),
						Recommendation: "Ограничьте прослушивание конкретным интерфейсом или добавьте аутентификацию/брандмауэр",
					})
				}
			}
		}
	})
	return issues
}
