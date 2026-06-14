package main

import (
	"flag"
	"io"
	"log"
	"os"

	"github.com/noedaka/go-config-parser/internal/config"
	"github.com/noedaka/go-config-parser/internal/parser"
	"github.com/noedaka/go-config-parser/internal/service"
	"github.com/noedaka/go-config-parser/internal/service/rules"
)

func main() {
	cfg := config.Init()

	var input []byte
	var err error

	if cfg.IsStdin {
		input, err = io.ReadAll(os.Stdin)
	} else {
		if flag.NArg() == 0 {
			log.Fatalf("Путь к файлу не указан")
		}

		input, err = os.ReadFile(flag.Arg(0))
	}

	if err != nil {
		log.Fatalf("Ошибка чтения: %v\n", err)
	}

	p := parser.Parser(parser.YamlJsonParser{})

	data, err := p.ParseConfig(input)
	if err != nil {
		log.Fatalf("Ошибка парсинга: %v\n", err)
	}

	rules := []service.Rule{
		rules.DebugLogRule{},
		rules.PlaintextPasswordRule{},
		rules.ZeroHostRule{},
		rules.TLSDisabledRule{},
		rules.NewWeakAlgorithmRule(),
	}

	var issues []service.Issue
	for _, rule := range rules {
		issues = append(issues, rule.Check(data)...)
	}

	for _, issue := range issues {
		log.Printf("[%s] %s\nРекомендация: %s\n\n", issue.Severity, issue.Message, issue.Recommendation)
	}

	if len(issues) > 0 && !cfg.IsSilent {
		log.Fatalf("Выходим с ошибкой, потому что есть проблемы в конфигурации.")
	}
	
	log.Print("В конфиге нет ошибок или флаг IsSilent активен.")
}
