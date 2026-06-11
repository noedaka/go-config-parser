package rules

import (
	"fmt"
	"strings"

	"github.com/noedaka/go-config-parser/internal/service"
)

type WeakAlgorithmRule struct {
	weak map[string]bool
}

func NewWeakAlgorithmRule() *WeakAlgorithmRule {
	return &WeakAlgorithmRule{
		weak: map[string]bool{
			"md5": true, "sha1": true, "rc4": true, "des": true, "3des": true,
			"md4": true, "ripemd160": true,
		},
	}
}

func (r *WeakAlgorithmRule) Name() string { return "weak-algorithm" }

func (r *WeakAlgorithmRule) Check(data any) []service.Issue {
	var issues []service.Issue
	algoKeys := []string{"algorithm", "algo", "digest-algorithm", "hash", "cipher", "encryption"}
	walk(data, func(key string, value any) {
		kl := strings.ToLower(key)
		for _, ak := range algoKeys {
			if kl == ak {
				if s, ok := value.(string); ok {
					if r.weak[strings.ToLower(s)] {
						issues = append(issues, service.Issue{
							Severity:       service.High,
							Message:        fmt.Sprintf("слабый алгоритм %s в поле '%s'", s, key),
							Recommendation: "Замените на стойкий алгоритм (SHA-256, AES-256 и т.п.)",
						})
					}
				}
			}
		}
	})
	return issues
}
