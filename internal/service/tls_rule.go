package service

import "strings"

type TLSDisabledRule struct{}

func (r TLSDisabledRule) Name() string { return "tls-disabled" }

func (r TLSDisabledRule) Check(data any) []Issue {
	var issues []Issue
	walk(data, func(key string, value any) {
		if !strings.EqualFold(key, "tls") && !strings.EqualFold(key, "ssl") {
			return
		}
		tlsMap, ok := value.(map[string]any)
		if !ok {
			return
		}
		// проверка на прямое отключение
		if v, exists := tlsMap["enabled"]; exists {
			if b, ok := v.(bool); ok && !b {
				issues = append(issues, Issue{
					Severity:       High,
					Message:        "TLS отключён (enabled: false)",
					Recommendation: "Включите TLS для всех внешних соединений",
				})
			}
		}
		// проверка небезопасных настроек проверки
		if v, exists := tlsMap["insecure_skip_verify"]; exists {
			if b, ok := v.(bool); ok && b {
				issues = append(issues, Issue{
					Severity:       High,
					Message:        "отключена проверка TLS-сертификата (insecure_skip_verify: true)",
					Recommendation: "Установите insecure_skip_verify в false и используйте доверенные сертификаты",
				})
			}
		}
		if v, exists := tlsMap["verify"]; exists {
			if b, ok := v.(bool); ok && !b {
				issues = append(issues, Issue{
					Severity:       High,
					Message:        "отключена верификация сертификата (verify: false)",
					Recommendation: "Включите проверку сертификатов",
				})
			}
		}
	})
	return issues
}
