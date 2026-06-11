package service

import (
	"fmt"
	"strings"
)

type PlaintextPasswordRule struct{}

func (r PlaintextPasswordRule) Name() string { return "plaintext-password" }

func (r PlaintextPasswordRule) Check(data any) []Issue {
	var issues []Issue
	suspiciousKeys := []string{"password", "passwd", "secret", "api_key", "token"}
	walk(data, func(key string, value any) {
		keyLower := strings.ToLower(key)
		for _, sk := range suspiciousKeys {
			if strings.Contains(keyLower, sk) {
				if s, ok := value.(string); ok && s != "" {
					// пропускаем хэшированные строки
					if strings.HasPrefix(s, "$") || strings.HasPrefix(s, "${") {
						continue
					}
					issues = append(issues, Issue{
						Severity:       High,
						Message:        fmt.Sprintf("пароль (поле '%s') хранится в открытом виде", key),
						Recommendation: "Используйте хэширование (bcrypt, argon2) или внешнее хранилище секретов",
					})
				}
			}
		}
	})
	return issues
}
